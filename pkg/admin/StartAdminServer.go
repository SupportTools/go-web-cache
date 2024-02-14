package admin

import (
	"log"
	"net/http"
	"strconv"

	"github.com/supporttools/go-web-cache/pkg/cache"
	"github.com/supporttools/go-web-cache/pkg/health"
)

// StartAdminServer initializes the admin server with the given cache manager.
func StartAdminServer(manager cache.CacheManager, port int) {
	cacheManager = manager

	// Health check endpoint
	http.HandleFunc("/healthz", health.HealthzHandler())
	http.HandleFunc("/readyz", health.ReadyzHandler())
	http.HandleFunc("/version", health.VersionHandler())
	http.HandleFunc("/api/items", handleListCacheItems)
	http.HandleFunc("/api/items/delete", handleDeleteCacheItem)

	// Convert port to string for ListenAndServe
	serverPort := strconv.Itoa(port)
	log.Printf("Admin server starting on port %d\n", port)
	if err := http.ListenAndServe(":"+serverPort, nil); err != nil {
		log.Fatalf("Admin server failed to start: %v", err)
	}
}
