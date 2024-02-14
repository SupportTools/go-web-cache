package cache

// NewInMemoryCache creates a new InMemoryCache with a specified maximum size and cache configuration.
func NewInMemoryCache(maxSize int, config CacheConfig) *InMemoryCache {
	return &InMemoryCache{
		items:       make(map[string]CacheItem),
		maxSize:     maxSize,
		keysOrder:   make([]string, 0),
		cacheConfig: config,
	}
}
