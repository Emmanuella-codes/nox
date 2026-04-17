package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	RedisURL      string
	JWTSecret     string
	Environment   string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port:        getEnv("PORT", "4006"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		RedisURL:    getEnv("REDIS_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
