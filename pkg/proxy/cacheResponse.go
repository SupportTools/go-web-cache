package proxy

import (
	"log"
	"net/http"

	"github.com/supporttools/go-web-cache/pkg/cache"
)

// cacheResponse determines whether the response should be cached and performs the caching if necessary.
func cacheResponse(t *Transport, cacheKey string, resp *http.Response, body []byte, req *http.Request) bool {
	// Log the beginning of the caching decision process
	log.Printf("Attempting to cache response for path: %s, Cache Key: %s", req.URL.Path, cacheKey)

	// Check and log the cache-control header
	cacheControlHeader := resp.Header.Get("Cache-Control")
	log.Printf("Cache-Control for %s: %s", req.URL.Path, cacheControlHeader)

	// Parse the Cache-Control header to determine if the response is cacheable
	cacheControl := cache.ParseCacheControl(cacheControlHeader)

	// You might have a custom function `ShouldCache` that determines if the response should be cached
	// based on the parsed Cache-Control directives, response status code, and other factors.
	if shouldCache := t.CacheManager.ShouldCache(resp, cacheControl); shouldCache {
		// If the response is deemed cacheable, proceed to cache it
		item := cache.CacheItem{
			ContentType:     resp.Header.Get("Content-Type"),
			Content:         body,
			CacheControl:    cache.ReconstructCacheControl(cacheControl),
			ContentEncoding: resp.Header.Get("Content-Encoding"),
			Path:            req.URL.Path,
		}

		// Set item expiration based on cache-control directives
		t.CacheManager.SetItemExpiration(&item, cacheControl)

		// Write the item to the cache with a default expiration if applicable
		t.CacheManager.WriteWithDefaultExpiration(cacheKey, item)

		// Log that the response has been cached
		log.Printf("Cached response for %s, Cache Key: %s", req.URL.Path, cacheKey)
		return true
	} else {
		// If the response should not be cached, log the reason
		log.Printf("Not caching response for %s due to Cache-Control or other policy: %s", req.URL.Path, cacheControlHeader)
		return false
	}
}
