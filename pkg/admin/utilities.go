package admin

import (
	"github.com/supporttools/go-web-cache/pkg/cache" // Ensure this path matches your project structure
)

// cacheManager holds the cache manager that admin functions will use to interact with the cache.
var cacheManager cache.CacheManager
