package cache

import (
	"net/http"
)

// ShouldCache decides whether the response should be cached based on Cache-Control directives
// and the content type, with a default caching policy for cacheable content types.
func (c *InMemoryCache) ShouldCache(resp *http.Response, cacheControl map[string]string) bool {
	contentType := resp.Header.Get("Content-Type")

	// Check if the content type is cacheable using this instance's configuration.
	isContentTypeCacheable := c.IsCacheableContentType(contentType)

	Logger.Debugf("ShouldCache: Checking if content type '%s' is cacheable.", contentType)
	Logger.Debugf("ShouldCache: Cacheable content types: %v", c.cacheConfig.CacheableMIMETypes)

	// Default to caching if the content type is cacheable and there's no "no-store" directive,
	// or if Cache-Control headers are absent.
	if len(cacheControl) == 0 && isContentTypeCacheable {
		Logger.Debugf("ShouldCache: Defaulting to cache for content type '%s'.", contentType)
		return resp.StatusCode == 200
	}

	// Do not cache if Cache-Control header includes no-store or private.
	if _, noStore := cacheControl["no-store"]; noStore {
		Logger.Debugf("ShouldCache: 'no-store' directive found. Not caching content type '%s'.", contentType)
		return false
	}
	if _, private := cacheControl["private"]; private {
		Logger.Debugf("ShouldCache: 'private' directive found. Not caching content type '%s'.", contentType)
		return false
	}

	// Cache if the response is successful (HTTP 200) and the content type is explicitly cacheable.
	cachable := resp.StatusCode == 200 && isContentTypeCacheable
	Logger.Debugf("ShouldCache: Final decision for content type '%s': %t", contentType, cachable)
	return cachable
}
