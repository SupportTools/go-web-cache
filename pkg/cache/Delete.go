package cache

// Delete removes an item from the cache.
func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.items[key]; exists {
		delete(c.items, key)
		// Remove the key from keysOrder
		for i, k := range c.keysOrder {
			if k == key {
				c.keysOrder = append(c.keysOrder[:i], c.keysOrder[i+1:]...)
				break
			}
		}
	}
}
