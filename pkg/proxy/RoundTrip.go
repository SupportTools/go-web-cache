package proxy

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/supporttools/go-web-cache/pkg/cache"
	"github.com/supporttools/go-web-cache/pkg/config"
	"github.com/supporttools/go-web-cache/pkg/metrics"
	"github.com/supporttools/go-web-cache/pkg/security"
)

// RoundTrip executes a single HTTP transaction, adding caching logic.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	metrics.IncrementTotalRequests()

	if req.Method != "GET" {
		return t.RoundTripper.RoundTrip(req)
	}

	modifiedReq := cloneRequestForClient(req)
	modifiedReq.URL.Scheme = config.CFG.BackendScheme
	modifiedReq.URL.Host = config.CFG.BackendServer
	modifiedReq.Host = config.CFG.BackendServer

	if security.HasWordPressLoginCookie(req) {
		return t.RoundTripper.RoundTrip(req)
	}

	// // Apply timeout through context
	// timeout := time.Duration(config.CFG.BackendTimeoutMs) * time.Millisecond
	// ctx, cancel := context.WithTimeout(req.Context(), timeout)
	// defer cancel() // Ensure the context is cancelled to prevent a context leak
	// modifiedReq = modifiedReq.WithContext(ctx)

	startTime := time.Now()
	cacheKey := cache.GetCacheKey(req)
	if item, found := t.CacheManager.Read(cacheKey); found && item.Expiration.After(time.Now()) {
		metrics.IncrementCacheHits()
		cacheHitStartTime := time.Now()
		cacheHitDuration := time.Since(cacheHitStartTime).Seconds()
		metrics.ObserveCacheHitResponseTime(cacheHitDuration)
		duration := time.Since(startTime).Seconds()
		metrics.ObserveTotalResponseTime(duration)
		return prepareCachedResponse(&item), nil
	}

	resp, err := t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	body, readErr := io.ReadAll(resp.Body)
	resp.Body.Close()
	if readErr != nil {
		return nil, readErr
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	metrics.IncrementCacheMisses()
	cacheMissStartTime := time.Now()

	go func() {
		if shouldCache := cacheResponse(t, cacheKey, resp, body, req); shouldCache {
		}
	}()

	cacheMissDuration := time.Since(cacheMissStartTime).Seconds()
	metrics.ObserveCacheMissResponseTime(cacheMissDuration)
	duration := time.Since(startTime).Seconds()
	metrics.ObserveTotalResponseTime(duration)

	return resp, nil
}

func cacheResponse(t *Transport, cacheKey string, resp *http.Response, body []byte, req *http.Request) bool {
	cacheControl := cache.ParseCacheControl(resp.Header.Get("Cache-Control"))
	if shouldCache := t.CacheManager.ShouldCache(resp, cacheControl); shouldCache {
		item := cache.CacheItem{
			ContentType:     resp.Header.Get("Content-Type"),
			Content:         body,
			CacheControl:    cache.ReconstructCacheControl(cacheControl),
			ContentEncoding: resp.Header.Get("Content-Encoding"),
			Path:            req.URL.Path,
		}
		t.CacheManager.SetItemExpiration(&item, cacheControl)
		t.CacheManager.WriteWithDefaultExpiration(cacheKey, item)
		return true
	}
	return false
}

func cloneRequestForClient(req *http.Request) *http.Request {
	// Deep copy the URL to ensure modifications don't affect the original request
	urlCopy := *req.URL
	clonedReq := req.WithContext(req.Context()) // Creates a shallow copy of the request
	clonedReq.URL = &urlCopy                    // Assign the deep copied URL

	// Explicitly copy the Host as WithContext does not do this
	clonedReq.Host = req.Host

	// No need to clear the RequestURI for client requests; it's ignored
	// Ensure other necessary headers or attributes are copied as needed

	return clonedReq
}
