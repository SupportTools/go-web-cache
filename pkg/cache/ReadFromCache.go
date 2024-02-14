package cache

import "time"

// ReadFromCache is a convenience function to access cached items.
// This could be a method of InMemoryCache if it needs to access the cache's internal state.
func (c *InMemoryCache) ReadFromCache(key string) (CacheItem, bool) {
	item, found := c.Read(key) // Use the receiver's Read method
	if !found || item.Expiration.Before(time.Now()) {
		return CacheItem{}, false
	}
	return item, true
}
