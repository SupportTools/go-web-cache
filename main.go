package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/supporttools/go-web-cache/pkg/admin"
	"github.com/supporttools/go-web-cache/pkg/cache"
	"github.com/supporttools/go-web-cache/pkg/config"
	"github.com/supporttools/go-web-cache/pkg/logging"
	"github.com/supporttools/go-web-cache/pkg/metrics"
	"github.com/supporttools/go-web-cache/pkg/proxy"
)

var logger = logging.SetupLogging()
var _ cache.CacheManager = (*cache.InMemoryCache)(nil)

func main() {
	logger.Println("Starting GO-Web-Cache")
	config.LoadConfiguration()
	if config.CFG.Debug {
		logger.Println("Debug mode enabled")
		logger.Println("Configuration:")
		logger.Printf("ConfigFile: %s", config.CFG.ConfigFile)
		logger.Printf("Debug: %t", config.CFG.Debug)
		logger.Printf("MetricsPort: %d", config.CFG.MetricsPort)
		logger.Printf("AdminPort: %d", config.CFG.AdminPort)
		logger.Printf("BackendHost: %s", config.CFG.BackendHost)
		logger.Printf("BackendScheme: %s", config.CFG.BackendScheme)
		logger.Printf("BackendPort: %d", config.CFG.BackendPort)
		logger.Printf("FrontendPort: %d", config.CFG.FrontendPort)
		logger.Printf("CacheMaxSize: %d", config.CFG.CacheMaxSize)
		logger.Println("CacheableMIMETypes:")
		for mimeType, cacheable := range config.CFG.CacheableMIMETypes {
			logger.Printf("  %s: %t", mimeType, cacheable)
		}
		logger.Println("NonCacheablePatterns:")
		for _, pattern := range config.CFG.NonCacheablePatterns {
			logger.Printf("  %s", pattern)
		}
	}

	cacheConfig := cache.CacheConfig{
		CacheableMIMETypes: config.CFG.CacheableMIMETypes,
	}
	cacheManager := cache.NewInMemoryCache(config.CFG.CacheMaxSize, cacheConfig)

	// Initialize and start the metrics and admin servers
	go metrics.StartMetricsServer(config.CFG.MetricsPort)
	go admin.StartAdminServer(cacheManager, config.CFG.AdminPort)

	// Setup reverse proxy and start the main server
	backendURL := fmt.Sprintf("%s://%s:%d", config.CFG.BackendScheme, config.CFG.BackendHost, config.CFG.BackendPort)
	proxyServer := proxy.NewReverseProxy(backendURL, cacheManager)
	http.HandleFunc("/", logging.LogRequests(proxyServer.ServeHTTP))

	log.Printf("Listening on port %d\n", config.CFG.FrontendPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.CFG.FrontendPort), nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
