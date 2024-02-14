package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/supporttools/go-web-cache/pkg/config"
	"github.com/supporttools/go-web-cache/pkg/logging"
)

// VersionInfo represents the structure of version information.
type VersionInfo struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	BuildTime string `json:"buildTime"`
}

var logger = logging.SetupLogging()
var version = "MISSING VERSION INFO"
var GitCommit = "MISSING GIT COMMIT"
var BuildTime = "MISSING BUILD TIME"

func HealthzHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	}
}

// ReadyzHandler checks if the application is ready and the backend server is up.
func ReadyzHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		backendURL := config.CFG.BackendServer + ":" + fmt.Sprint(config.CFG.BackendPort) + fmt.Sprint(config.CFG.BackendHealthCheck)

		// Perform a simple GET request to the backend's health check endpoint
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(backendURL)
		if err != nil || resp.StatusCode != 200 {
			logger.Error("Backend server check failed", err)
			http.Error(w, "Backend server is not ready", http.StatusServiceUnavailable)
			return
		}

		if config.CFG.Debug {
			logger.Info("ReadyzHandler: Backend server is up")
		}
		fmt.Fprintf(w, "ok")
	}
}

// VersionHandler returns version information as JSON.
func VersionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("VersionHandler")

		versionInfo := VersionInfo{
			Version:   version,
			GitCommit: GitCommit,
			BuildTime: BuildTime,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(versionInfo); err != nil {
			logger.Error("Failed to encode version info to JSON", err)
			http.Error(w, "Failed to encode version info", http.StatusInternalServerError)
		}
	}
}
