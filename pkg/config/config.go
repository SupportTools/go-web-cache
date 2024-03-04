package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

// AppConfig structure combines both file-based and environment-based configurations.
type AppConfig struct {
	ConfigFile           string          `json:"configFile"`
	Debug                bool            `json:"debug"`
	MetricsPort          int             `json:"metricsPort"`
	AdminPort            int             `json:"adminPort"`
	BackendServer        string          `json:"backendServer"`
	BackendPort          int             `json:"backendPort"`
	BackendHealthCheck   string          `json:"backendHealthCheck"`
	BackendTimeoutMs     int             `json:"backendTimeoutMs"`
	FrontendPort         int             `json:"frontendPort"`
	CacheMaxSize         int             `json:"cacheMaxSize"`
	CacheableMIMETypes   map[string]bool `json:"cacheableMIMETypes"`
	NonCacheablePatterns []string        `json:"nonCacheablePatterns"`
}

var CFG AppConfig

// LoadConfiguration loads configuration from a JSON file and environment variables.
func LoadConfiguration() {
	filename := getEnvOrDefault("CONFIG_FILE", "./config.json")
	if filename != "" {
		err := loadConfigFile(filename)
		if err != nil {
			log.Printf("Warning: Failed to load config file %s: %v", filename, err)
		}
	}

	overrideConfigWithEnv()
}

func loadConfigFile(filename string) error {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(fileContent, &CFG)
}

func overrideConfigWithEnv() {
	CFG.ConfigFile = getEnvOrDefault("CONFIG_FILE", CFG.ConfigFile)
	CFG.Debug = parseEnvBool("DEBUG", CFG.Debug)
	CFG.MetricsPort = parseEnvInt("METRICS_PORT", CFG.MetricsPort)
	CFG.AdminPort = parseEnvInt("ADMIN_PORT", CFG.AdminPort)
	CFG.BackendServer = getEnvOrDefault("BACKEND_SERVER", CFG.BackendServer)
	CFG.BackendPort = parseEnvInt("BACKEND_PORT", CFG.BackendPort)
	CFG.BackendHealthCheck = getEnvOrDefault("BACKEND_HEALTH_CHECK", CFG.BackendHealthCheck)
	CFG.BackendTimeoutMs = parseEnvInt("BACKEND_TIMEOUT", CFG.BackendTimeoutMs)
	CFG.FrontendPort = parseEnvInt("FRONTEND_PORT", CFG.FrontendPort)
	CFG.CacheMaxSize = parseEnvInt("CACHE_MAX_SIZE", CFG.CacheMaxSize)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func parseEnvInt(key string, defaultValue int) int {
	value := getEnvOrDefault(key, "")
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Error parsing %s as int: %v. Using default value: %d", key, err, defaultValue)
		return defaultValue
	}
	return intValue
}

func parseEnvBool(key string, defaultValue bool) bool {
	value := getEnvOrDefault(key, "")
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Error parsing %s as bool: %v. Using default value: %t", key, err, defaultValue)
		return defaultValue
	}
	return boolValue
}
