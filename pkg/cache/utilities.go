// cache/cache.go

package cache

import (
	"net/http"
	"sync"
	"time"

	"github.com/supporttools/go-web-cache/pkg/logging"
)

var Logger = logging.SetupLogging()

// CacheConfig holds configuration for the cache, including cacheable MIME types.
type CacheConfig struct {
	CacheableMIMETypes map[string]bool
	DefaultExpiration  time.Duration
}

// CacheItem represents an item in the cache.
type CacheItem struct {
	Path            string
	ContentType     string
	Content         []byte
	ContentEncoding string
	Expiration      time.Time
	CacheControl    string
}

// CacheManager defines the interface for interacting with the cache.
type CacheManager interface {
	Read(key string) (CacheItem, bool)
	Write(key string, item CacheItem)
	Delete(key string)
	List() []string
	ShouldCache(resp *http.Response, cacheControl map[string]string) bool
	SetItemExpiration(item *CacheItem, cacheControl map[string]string)
	IsCacheableContentType(contentType string) bool
	WriteWithExpiration(key string, item CacheItem, cacheControl map[string]string)
	WriteWithDefaultExpiration(key string, item CacheItem)
}

// InMemoryCache provides a simple in-memory cache implementation of CacheManager.
type InMemoryCache struct {
	items       map[string]CacheItem
	mu          sync.RWMutex
	maxSize     int
	keysOrder   []string
	cacheConfig CacheConfig
}
