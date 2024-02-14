package cache

import "time"

// Write adds or updates an item in the cache with the default expiration if not set.
func (c *InMemoryCache) WriteWithDefaultExpiration(key string, item CacheItem) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Apply default expiration if not set
	if item.Expiration.IsZero() {
		// Check for a default expiration in the configuration
		if c.cacheConfig.DefaultExpiration > 0 {
			item.Expiration = time.Now().Add(c.cacheConfig.DefaultExpiration)
		} else {
			// Apply a hard-coded default if no configuration is present
			item.Expiration = time.Now().Add(60 * time.Minute) // Example: 60 minutes
		}
	}

	// Proceed to cache the item
	c.items[key] = item
	c.keysOrder = append(c.keysOrder, key)
}
