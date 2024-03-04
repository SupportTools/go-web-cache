package proxy

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/supporttools/go-web-cache/pkg/cache"
	"github.com/supporttools/go-web-cache/pkg/metrics"
	"github.com/supporttools/go-web-cache/pkg/security"
)

// RoundTrip executes a single HTTP transaction, adding caching logic.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Printf("RoundTrip called for URL: %s", req.URL.String())
	metrics.IncrementTotalRequests()

	// if req.Method != "GET" {
	// 	log.Println("Non-GET request, bypassing cache and forwarding directly")
	// 	return t.RoundTripper.RoundTrip(req)
	// }

	//modifiedReq := cloneRequestForClient(req)
	//log.Printf("Modified request for caching: URL Scheme: %s, Host: %s", modifiedReq.URL.Scheme, modifiedReq.Host)

	if security.HasWordPressLoginCookie(req) {
		log.Println("Request has WordPress login cookie, bypassing cache")
		return t.RoundTripper.RoundTrip(req)
	}

	startTime := time.Now()
	cacheKey := cache.GetCacheKey(req)
	log.Printf("Cache key generated: %s", cacheKey)
	if item, found := t.CacheManager.Read(cacheKey); found && item.Expiration.After(time.Now()) {
		log.Println("Cache hit")
		metrics.IncrementCacheHits()
		cacheHitStartTime := time.Now()
		cacheHitDuration := time.Since(cacheHitStartTime).Seconds()
		metrics.ObserveCacheHitResponseTime(cacheHitDuration)
		duration := time.Since(startTime).Seconds()
		metrics.ObserveTotalResponseTime(duration)
		return prepareCachedResponse(&item), nil
	}

	log.Println("Cache miss, forwarding request to backend")
	resp, err := t.RoundTripper.RoundTrip(req)
	if err != nil {
		log.Printf("Error forwarding request to backend: %v", err)
		return nil, err
	}

	body, readErr := io.ReadAll(resp.Body)
	resp.Body.Close()
	if readErr != nil {
		log.Printf("Error reading response body: %v", readErr)
		return nil, readErr
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	metrics.IncrementCacheMisses()

	cacheMissStartTime := time.Now()
	go func() {
		if shouldCache := cacheResponse(t, cacheKey, resp, body, req); shouldCache {
			log.Printf("Response for %s cached successfully", req.URL.Path)
		} else {
			log.Printf("Response for %s not cached", req.URL.Path)
		}
	}()

	cacheMissDuration := time.Since(cacheMissStartTime).Seconds()
	metrics.ObserveCacheMissResponseTime(cacheMissDuration)
	duration := time.Since(startTime).Seconds()
	metrics.ObserveTotalResponseTime(duration)

	return resp, nil
}

func cacheResponse(t *Transport, cacheKey string, resp *http.Response, body []byte, req *http.Request) bool {
	// Log for debugging
	log.Printf("Attempting to cache response for path: %s, Cache Key: %s", req.URL.Path, cacheKey)

	// Check and log the cache-control header
	cacheControlHeader := resp.Header.Get("Cache-Control")
	log.Printf("Cache-Control for %s: %s", req.URL.Path, cacheControlHeader)

	cacheControl := cache.ParseCacheControl(cacheControlHeader)
	if shouldCache := t.CacheManager.ShouldCache(resp, cacheControl); shouldCache {
		// Additional logging to confirm caching decision
		log.Printf("Caching enabled for: %s", req.URL.Path)
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
	} else {
		// Log why caching was not performed
		log.Printf("Response for %s not cached, Cache-Control: %s", req.URL.Path, cacheControlHeader)
	}
	return false
}

// func cloneRequestForClient(req *http.Request) *http.Request {
// 	urlCopy := *req.URL
// 	clonedReq := req.WithContext(req.Context())
// 	clonedReq.URL = &urlCopy
// 	clonedReq.Host = req.Host
// 	log.Printf("Request cloned for modification: %s", clonedReq.URL.String())
// 	return clonedReq
// }
