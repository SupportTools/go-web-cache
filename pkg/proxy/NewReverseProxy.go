package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/supporttools/go-web-cache/pkg/cache"
)

// NewReverseProxy creates and configures a reverse proxy.
func NewReverseProxy(targetURL string, cacheManager cache.CacheManager) *httputil.ReverseProxy {
	url, _ := url.Parse(targetURL)

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = &Transport{http.DefaultTransport, cacheManager}

	return proxy
}
