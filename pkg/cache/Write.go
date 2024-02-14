package cache

// Write adds or updates an item in the cache.
func (c *InMemoryCache) Write(key string, item CacheItem) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if the item already exists to avoid unnecessary eviction
	if _, exists := c.items[key]; !exists {
		// Evict the oldest item if we're at max capacity
		if len(c.items) >= c.maxSize {
			oldestKey := c.keysOrder[0]
			delete(c.items, oldestKey)
			c.keysOrder = c.keysOrder[1:] // Remove the oldest item key
		}
		c.keysOrder = append(c.keysOrder, key) // Add new item key as the newest
	}

	c.items[key] = item
}
