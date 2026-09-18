package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL     string
	PublicAPIKey    string
	AdminAPIKey     string
	Port            string
	AllowedOrigin   string
	AllowCredentials bool
	DevMode         bool
}

func Load() (Config, error) {
	devMode, _ := strconv.ParseBool(os.Getenv("DEV_MODE"))
	allowCredentials, _ := strconv.ParseBool(os.Getenv("ALLOW_CREDENTIALS"))
	cfg := Config{
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		PublicAPIKey:    os.Getenv("TELEMETRY_PUBLIC_API_KEY"),
		AdminAPIKey:     os.Getenv("TELEMETRY_ADMIN_API_KEY"),
		Port:            os.Getenv("PORT"),
		AllowedOrigin:   os.Getenv("ALLOWED_ORIGIN"),
		AllowCredentials: allowCredentials,
		DevMode:         devMode,
	}
	if cfg.DatabaseURL == "" {
		return cfg, errors.New("DATABASE_URL is required")
	}
	if !cfg.DevMode {
		if cfg.PublicAPIKey == "" {
			return cfg, errors.New("TELEMETRY_PUBLIC_API_KEY is required")
		}
		if cfg.AdminAPIKey == "" {
			return cfg, errors.New("TELEMETRY_ADMIN_API_KEY is required")
		}
	}
	if cfg.Port == "" {
		cfg.Port = "9999"
	}
	if cfg.AllowedOrigin == "" {
		cfg.AllowedOrigin = "*"
	}
	return cfg, nil
}
