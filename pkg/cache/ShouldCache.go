package cache

import (
	"net/http"
	"strconv"
	"strings"
)

// ShouldCache decides whether the response should be cached based on Cache-Control directives
// and the content type, with a default caching policy for cacheable content types.
func (c *InMemoryCache) ShouldCache(resp *http.Response, cacheControl map[string]string) bool {
	contentType := resp.Header.Get("Content-Type")

	// Extract the base content type without parameters (e.g., charset).
	baseContentType := strings.Split(contentType, ";")[0]

	// Check if the content type is cacheable using this instance's configuration.
	isContentTypeCacheable := c.IsCacheableContentType(baseContentType)

	Logger.Debugf("ShouldCache: Checking if content type '%s' (base: '%s') is cacheable.", contentType, baseContentType)
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

	// Cache if max-age directive is present and has a positive value.
	if maxAge, exists := cacheControl["max-age"]; exists {
		maxAgeValue, err := strconv.Atoi(maxAge)
		if err == nil && maxAgeValue > 0 {
			Logger.Debugf("ShouldCache: 'max-age=%d' directive found. Caching content type '%s'.", maxAgeValue, contentType)
			return true
		}
	}

	// Cache if the public directive is present.
	if _, public := cacheControl["public"]; public {
		Logger.Debugf("ShouldCache: 'public' directive found. Caching content type '%s'.", contentType)
		return true
	}

	// Cache if the response is successful (HTTP 200) and the content type is explicitly cacheable.
	cacheable := resp.StatusCode == 200 && isContentTypeCacheable
	Logger.Debugf("ShouldCache: Final decision for content type '%s': %t", contentType, cacheable)
	return cacheable
}
