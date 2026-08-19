// Package config loads runtime configuration from the environment with safe
// defaults, so the service runs with zero configuration in development while
// remaining fully configurable in production (Principle of Least Astonishment).
package config

import (
	"os"
	"strings"
)

// Config is the resolved runtime configuration.
type Config struct {
	// Port is the TCP port the HTTP server listens on.
	Port string
	// AllowedOrigins is the list of origins permitted by CORS. A single "*"
	// allows any origin (development default).
	AllowedOrigins []string
}

// Load reads configuration from the environment, applying defaults for any
// unset value.
//
//	PORT           TCP port to listen on            (default "8080")
//	CORS_ORIGIN    comma-separated allowed origins  (default "*")
func Load() Config {
	return Config{
		Port:           getenv("PORT", "8080"),
		AllowedOrigins: splitAndTrim(getenv("CORS_ORIGIN", "*")),
	}
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return v
	}
	return fallback
}

func splitAndTrim(csv string) []string {
	parts := strings.Split(csv, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
