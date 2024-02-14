package cache

import "time"

// WriteToCache stores an item in the cache with an optional expiration based on cacheControl directives.
func (c *InMemoryCache) WriteToCache(key string, item CacheItem, cacheControl map[string]string) {
	// Check for max-age directive in cacheControl to set item expiration
	if maxAge, ok := cacheControl["max-age"]; ok {
		maxAgeDuration, err := time.ParseDuration(maxAge + "s")
		if err == nil {
			item.Expiration = time.Now().Add(maxAgeDuration)
		}
	}

	// Use the Write method of InMemoryCache to store the item
	c.Write(key, item)
}
