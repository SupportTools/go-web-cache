package cache

// IsCacheableContentType checks if the content type is cacheable based on the configuration.
func (c *InMemoryCache) IsCacheableContentType(contentType string) bool {
	_, cacheable := c.cacheConfig.CacheableMIMETypes[contentType]
	return cacheable
}
