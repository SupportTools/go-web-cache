package cache

import (
	"net/http"
	"sort"
	"strings"
)

func GetCacheKey(req *http.Request) string {
	var cacheKey strings.Builder
	cacheKey.WriteString(req.URL.Path)

	// Include query parameters in cache key
	queryParams := req.URL.Query()
	if len(queryParams) > 0 {
		// Sort query parameters to ensure consistent ordering
		var keys []string
		for k := range queryParams {
			keys = append(keys, k)
		}
		sort.Strings(keys) // Ensure the order is consistent
		cacheKey.WriteString("?")
		for i, k := range keys {
			if i > 0 {
				cacheKey.WriteString("&")
			}
			cacheKey.WriteString(k)
			cacheKey.WriteString("=")
			cacheKey.WriteString(queryParams.Get(k))
		}
	}

	// Incorporate Vary header values if present and relevant
	// Assuming `varyHeaders` are obtained from a previous response or configuration
	varyHeaders := req.Header.Get("Vary")
	if varyHeaders != "" {
		cacheKey.WriteString("#")
		for _, headerName := range strings.Split(varyHeaders, ",") {
			headerValue := req.Header.Get(strings.TrimSpace(headerName))
			if headerValue != "" {
				cacheKey.WriteString(headerName)
				cacheKey.WriteString("=")
				cacheKey.WriteString(headerValue)
				cacheKey.WriteString(";")
			}
		}
	}

	return cacheKey.String()
}
