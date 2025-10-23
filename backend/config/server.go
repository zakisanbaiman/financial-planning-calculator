package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// ServerConfig holds server configuration
type ServerConfig struct {
	Port                string
	Debug               bool
	AllowedOrigins      []string
	RateLimitRPS        int
	RequestTimeout      time.Duration
	MaxRequestSize      string
	EnableGzip          bool
	GzipLevel           int
	LogFormat           string
	EnableSecureHeaders bool
}

// LoadServerConfig loads server configuration from environment variables
func LoadServerConfig() *ServerConfig {
	config := &ServerConfig{
		Port:                getEnv("PORT", "8080"),
		Debug:               getEnvBool("DEBUG", false),
		AllowedOrigins:      getEnvSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:3001", "https://localhost:3000", "https://localhost:3001"}),
		RateLimitRPS:        getEnvInt("RATE_LIMIT_RPS", 100),
		RequestTimeout:      getEnvDuration("REQUEST_TIMEOUT", 30*time.Second),
		MaxRequestSize:      getEnv("MAX_REQUEST_SIZE", "10M"),
		EnableGzip:          getEnvBool("ENABLE_GZIP", true),
		GzipLevel:           getEnvInt("GZIP_LEVEL", 5),
		LogFormat:           getEnv("LOG_FORMAT", "${time_rfc3339} ${method} ${uri} ${status} ${latency_human} ${bytes_in}B/${bytes_out}B ${error}\n"),
		EnableSecureHeaders: getEnvBool("ENABLE_SECURE_HEADERS", true),
	}

	return config
}

// Helper functions for environment variable parsing

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
