package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL   string
	APIKey        string
	Port          string
	AllowedOrigin string
}

func Load() (Config, error) {
	cfg := Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		APIKey:        os.Getenv("TELEMETRY_API_KEY"),
		Port:          os.Getenv("PORT"),
		AllowedOrigin: os.Getenv("ALLOWED_ORIGIN"),
	}
	if cfg.DatabaseURL == "" {
		return cfg, errors.New("DATABASE_URL is required")
	}
	if cfg.APIKey == "" {
		return cfg, errors.New("TELEMETRY_API_KEY is required")
	}
	if cfg.Port == "" {
		cfg.Port = "9999"
	}
	if cfg.AllowedOrigin == "" {
		cfg.AllowedOrigin = "*"
	}
	return cfg, nil
}
