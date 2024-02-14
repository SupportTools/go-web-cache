package cache

// WriteWithExpiration stores an item in the cache with an expiration time based on cacheControl.
func (c *InMemoryCache) WriteWithExpiration(key string, item CacheItem, cacheControl map[string]string) {
	c.SetItemExpiration(&item, cacheControl)
	c.Write(key, item)
}
