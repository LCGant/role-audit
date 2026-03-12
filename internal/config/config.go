package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HTTPAddr      string
	InternalToken string
	MetricsToken  string
	LogLevel      string
	LogFile       string
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:      getenv("AUDIT_HTTP_ADDR", ":8080"),
		InternalToken: strings.TrimSpace(os.Getenv("AUDIT_INTERNAL_TOKEN")),
		MetricsToken:  strings.TrimSpace(os.Getenv("AUDIT_METRICS_TOKEN")),
		LogLevel:      getenv("AUDIT_LOG_LEVEL", "info"),
		LogFile:       strings.TrimSpace(os.Getenv("AUDIT_LOG_FILE")),
	}
	if cfg.InternalToken == "" {
		return Config{}, errors.New("AUDIT_INTERNAL_TOKEN is required")
	}
	if cfg.LogFile == "" {
		return Config{}, errors.New("AUDIT_LOG_FILE is required")
	}
	return cfg, nil
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}
