package proxy

import (
	"bytes"
	"io"
	"net/http"

	"github.com/supporttools/go-web-cache/pkg/cache"
)

// prepareCachedResponse constructs an http.Response from a cached item.
func prepareCachedResponse(item *cache.CacheItem) *http.Response {
	headers := http.Header{
		"Content-Type":  []string{item.ContentType},
		"Cache-Control": []string{item.CacheControl},
	}
	if item.ContentEncoding != "" {
		headers.Set("Content-Encoding", item.ContentEncoding)
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(item.Content)),
		Header:     headers,
	}
}
