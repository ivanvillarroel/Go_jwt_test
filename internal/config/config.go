package config

import "os"

type Config struct {
	Port         string
	JWTSecret    string
	AuthUser     string
	AuthPassword string
}

func Load() Config {
	return Config{
		Port:         valueOrDefault("PORT", "8080"),
		JWTSecret:    valueOrDefault("JWT_SECRET", "local-development-secret"),
		AuthUser:     valueOrDefault("AUTH_USER", "api-user"),
		AuthPassword: valueOrDefault("AUTH_PASSWORD", "password"),
	}
}

func valueOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
