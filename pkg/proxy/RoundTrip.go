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
	// Increment the total number of requests.
	metrics.IncrementTotalRequests()

	// Skip caching for non-GET requests.
	if req.Method != "GET" {
		Logger.Println("Bypassing cache for non-GET request")
		return t.RoundTripper.RoundTrip(req)
	}

	// Setting Timeout for the request
	timeout := time.Duration(config.CFG.BackendTimeoutMs) * time.Millisecond
	client := &http.Client{
		Timeout:   timeout,
		Transport: t.RoundTripper,
	}

	// Skip caching for WordPress login cookies.
	if security.HasWordPressLoginCookie(req) {
		Logger.Println("Bypassing cache for logged-in WordPress user")
		//return t.RoundTripper.RoundTrip(req)
		return client.Do(req)
	}

	// Start timer to track overall response time
	startTime := time.Now()

	// Attempt to serve the request from cache.
	cacheKey := cache.GetCacheKey(req)
	if item, found := t.CacheManager.Read(cacheKey); found && item.Expiration.After(time.Now()) {
		// Cache hit
		metrics.IncrementCacheHits()

		// Calculate cache hit response time
		cacheHitStartTime := time.Now()

		Logger.Printf("Cache hit for: %s", cacheKey)
		if config.CFG.Debug {
			Logger.Printf("Debug: Serving %s from cache", req.URL.Path)
		}

		// Calculate cache hit response time duration
		cacheHitDuration := time.Since(cacheHitStartTime).Seconds()
		metrics.ObserveCacheHitResponseTime(cacheHitDuration)

		// Calculate overall response time duration
		duration := time.Since(startTime).Seconds()
		metrics.ObserveTotalResponseTime(duration)

		return prepareCachedResponse(&item), nil
	}

	// Cache miss or expired item, perform the request.
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	Logger.Debugf("Received response for %s with status code %d", req.URL.Path, resp.StatusCode)

	// Clone the response body for both caching and responding to the client.
	body, readErr := io.ReadAll(resp.Body)
	resp.Body.Close() // Close the original body.
	if readErr != nil {
		return nil, readErr
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	// Increment cache misses counter
	metrics.IncrementCacheMisses()

	// Calculate cache miss response time
	cacheMissStartTime := time.Now()

	// Response received, send it to the client.
	go func(t *Transport, cacheKey string, resp *http.Response, body []byte, req *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				Logger.Printf("Recovering from panic during cache write: %v", r)
			}
		}()
		if shouldCache := cacheResponse(t, cacheKey, resp, body, req); shouldCache {
			Logger.Debugf("Caching response for %s", cacheKey)
		}
	}(t, cacheKey, resp, body, req)

	// Calculate cache miss response time duration
	cacheMissDuration := time.Since(cacheMissStartTime).Seconds()
	metrics.ObserveCacheMissResponseTime(cacheMissDuration)

	// Calculate overall response time duration
	duration := time.Since(startTime).Seconds()
	metrics.ObserveTotalResponseTime(duration)

	return resp, nil
}

// cacheResponse caches the response if necessary.
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
	Logger.Debugf("Not caching response for: %s", cacheKey)
	return false
}
