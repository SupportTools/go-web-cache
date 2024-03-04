package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/supporttools/go-web-cache/pkg/cache"
	"github.com/supporttools/go-web-cache/pkg/config"
)

// NewReverseProxy creates and configures a reverse proxy.
func NewReverseProxy(targetURL string, cacheManager cache.CacheManager) *httputil.ReverseProxy {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Error parsing target URL '%s': %v", targetURL, err)
	}

	log.Printf("Creating reverse proxy to %s", parsedURL.String())

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	proxy.Transport = &Transport{
		RoundTripper: http.DefaultTransport,
		CacheManager: cacheManager,
	}

	// Modify proxy.Director for additional debugging and request customization
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		log.Printf("Proxying request for %s", req.URL.String())

		// Call the original director to preserve default behavior
		originalDirector(req)

		// Log the request URL to debug
		if config.CFG.Debug {
			log.Printf("Modified request URL: %s", req.URL.String())
		}

		// Ensure the request host is correctly set
		req.URL.Scheme = "http" // or "https", as appropriate
		req.URL.Host = parsedURL.Host
		req.Host = parsedURL.Host

		// Additional logging to confirm the host and scheme are correctly set
		if config.CFG.Debug {
			log.Printf("Request URL after modification: %s", req.URL.String())
			log.Printf("Request Host after modification: %s", req.Host)
			// Debugging: Log the request headers after modification by the director
			log.Printf("Request headers after director modifications:")
			for name, values := range req.Header {
				for _, value := range values {
					log.Printf("%s: %s", name, value)
				}
			}
		}

	}

	// Optionally, customize the Transport to log errors or responses
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Error handling request for %s: %v", r.URL.String(), err)
	}

	return proxy
}
