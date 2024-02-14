package cache

import "time"

// SetItemExpiration adjusts the expiration of a cache item based on Cache-Control directives.
func (c *InMemoryCache) SetItemExpiration(item *CacheItem, cacheControl map[string]string) {
	maxAgeStr, exists := cacheControl["max-age"]
	if exists {
		maxAge, err := time.ParseDuration(maxAgeStr + "s")
		if err == nil {
			item.Expiration = time.Now().Add(maxAge)
		}
	}
}
