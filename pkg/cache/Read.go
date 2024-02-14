package cache

// Read retrieves an item from the cache.
func (c *InMemoryCache) Read(key string) (CacheItem, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if found {
		Logger.Debugf("Cache hit for key: %s, Expiration: %s", key, item.Expiration)
	} else {
		Logger.Debugf("Cache miss for key: %s", key)
	}
	return item, found
}
