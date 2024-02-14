package metrics

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/supporttools/go-web-cache/pkg/logging"
)

var logger = logging.SetupLogging()

var (
	totalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "go_web_requests_total",
		Help: "Total number of requests.",
	})
	cacheHits = promauto.NewCounter(prometheus.CounterOpts{
		Name: "go_web_cache_hits_total",
		Help: "Total number of cache hits.",
	})
	cacheMisses = promauto.NewCounter(prometheus.CounterOpts{
		Name: "go_web_cache_misses_total",
		Help: "Total number of cache misses.",
	})
	cacheItemCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "go_web_cache_items",
		Help: "Current number of items in the cache.",
	})
	backendResponseTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "backend_response_time_seconds",
		Help:    "Histogram of response times from the backend.",
		Buckets: prometheus.DefBuckets,
	})
	cacheHitResponseTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "cache_hit_response_time_seconds",
		Help:    "Histogram of response times for cache hits.",
		Buckets: prometheus.DefBuckets,
	})
	cacheMissResponseTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "cache_miss_response_time_seconds",
		Help:    "Histogram of response times for cache misses.",
		Buckets: prometheus.DefBuckets,
	})
	totalResponseTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "total_response_time_seconds",
		Help:    "Histogram of overall response times.",
		Buckets: prometheus.DefBuckets,
	})
)

func init() {
	// Register the histograms with Prometheus.
	prometheus.MustRegister(cacheHitResponseTime)
	prometheus.MustRegister(cacheMissResponseTime)
	prometheus.MustRegister(totalResponseTime)
}

// StartMetricsServer starts the Prometheus metrics server on the specified port.
func StartMetricsServer(port int) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	// Convert port to string for ListenAndServe
	serverPort := strconv.Itoa(port)
	logger.Printf("Metrics server starting on port %d\n", port) // Updated to use %d for integer formatting
	if err := http.ListenAndServe(":"+serverPort, mux); err != nil {
		logger.Fatalf("Metrics server failed to start: %v", err)
	}
}

// IncrementCacheHits increments the cache hits counter.
func IncrementCacheHits() {
	cacheHits.Inc()
}

// IncrementCacheMisses increments the cache misses counter.
func IncrementCacheMisses() {
	cacheMisses.Inc()
}

// IncrementTotalRequests increments the total requests counter.
func IncrementTotalRequests() {
	totalRequests.Inc()
}

// SetCacheItemCount sets the current number of items in the cache.
func SetCacheItemCount(count float64) {
	cacheItemCount.Set(count)
}

// ObserveBackendResponseTime records the response time from the backend.
func ObserveBackendResponseTime(duration float64) {
	backendResponseTime.Observe(duration)
}

// ObserveCacheHitResponseTime records the response time for cache hits.
func ObserveCacheHitResponseTime(duration float64) {
	cacheHitResponseTime.Observe(duration)
}

// ObserveCacheMissResponseTime records the response time for cache misses.
func ObserveCacheMissResponseTime(duration float64) {
	cacheMissResponseTime.Observe(duration)
}

// ObserveTotalResponseTime records the overall response time.
func ObserveTotalResponseTime(duration float64) {
	backendResponseTime.Observe(duration)
}
