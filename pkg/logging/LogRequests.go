package logging

import (
	"fmt"
	"net/http"
	"time"

	"github.com/supporttools/go-web-cache/pkg/config"
)

// LogRequests wraps the HTTP handler to log requests
func LogRequests(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the request details
		if r.URL.Path == "/healthz" || r.URL.Path == "/readyz" || r.URL.Path == "/version" {
			next(w, r) // Call the next handler
			return
		}
		if config.CFG.Debug {
			fmt.Printf("Received request %s %s at %s\n", r.Method, r.URL.Path, time.Now().Format(time.RFC3339))
		}

		next(w, r) // Call the next handler
	}
}
