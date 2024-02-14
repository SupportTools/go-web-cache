package proxy

import (
	"net/http"

	"github.com/supporttools/go-web-cache/pkg/cache"
	"github.com/supporttools/go-web-cache/pkg/logging"
)

var Logger = logging.SetupLogging()

// Transport extends the http.RoundTripper interface with caching logic.
type Transport struct {
	http.RoundTripper
	CacheManager cache.CacheManager
}
