// Package config holds runtime configuration resolved from flags
// and environment variables. It's deliberately a plain struct — the
// exporter has a small surface area and we don't need a full CLI
// framework.
package config

import (
	"os"
	"strings"
	"time"
)

const envPrefix = "HY_SMI_EXPORTER_"

// Config is the resolved runtime configuration.
type Config struct {
	ListenAddr  string
	MetricsPath string
	LogLevel    string
	BinaryPath  string
	Timeout     time.Duration
}

// Default returns the zero-config with sensible defaults applied.
func Default() Config {
	return Config{
		ListenAddr:  ":9584",
		MetricsPath: "/metrics",
		LogLevel:    "info",
		BinaryPath:  "hy-smi",
		Timeout:     5 * time.Second,
	}
}

// FromEnv returns Default() overridden by HY_SMI_EXPORTER_* env vars.
func FromEnv() Config {
	c := Default()
	c.ListenAddr = envOr(envPrefix+"LISTEN_ADDR", c.ListenAddr)
	c.MetricsPath = envOr(envPrefix+"METRICS_PATH", c.MetricsPath)
	c.LogLevel = envOr(envPrefix+"LOG_LEVEL", c.LogLevel)
	c.BinaryPath = envOr(envPrefix+"BINARY_PATH", c.BinaryPath)
	if v, ok := os.LookupEnv(envPrefix + "TIMEOUT"); ok && v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			c.Timeout = d
		}
	}
	return c
}

func envOr(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

// LogLevelValue is a convenience for callers that want to validate
// the level at parse time.
func (c Config) LogLevelValue() string {
	return strings.ToLower(strings.TrimSpace(c.LogLevel))
}
