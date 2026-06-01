package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                  string
	DatabaseURL           string
	RedisURL              string
	JWTAccessSecret       string
	JWTRefreshSecret      string
	JWTIssuer             string
	JWTAudience           string
	JWTAccessTTL          time.Duration
	JWTRefreshTTL         time.Duration
	EmailOTPTTL           time.Duration
	BrevoAPIKey           string
	BrevoBaseURL          string
	MailFromEmail         string
	MailFromName          string
	MediaUploadBaseURL    string
	MediaPublicBaseURL    string
	MediaProcessingSecret string
	GhostPersonaSecret    string
	Environment           string
}

func Load() (*Config, error) {
	for _, f := range []string{".env", "../.env"} {
		if err := godotenv.Load(f); err == nil {
			break
		}
	}

	cfg := &Config{
		Port:                  getEnv("PORT", "4006"),
		DatabaseURL:           getEnv("DATABASE_URL", ""),
		RedisURL:              getEnv("REDIS_URL", ""),
		JWTAccessSecret:       getEnv("JWT_ACCESS_SECRET", getEnv("JWT_SECRET", "")),
		JWTRefreshSecret:      getEnv("JWT_REFRESH_SECRET", ""),
		JWTIssuer:             getEnv("JWT_ISSUER", "nox-api"),
		JWTAudience:           getEnv("JWT_AUDIENCE", "nox-client"),
		JWTAccessTTL:          getDurationEnv("JWT_ACCESS_TTL", 15*time.Minute),
		JWTRefreshTTL:         getDurationEnv("JWT_REFRESH_TTL", 720*time.Hour),
		EmailOTPTTL:           getDurationEnv("EMAIL_OTP_TTL", 10*time.Minute),
		BrevoAPIKey:           getEnv("BREVO_API_KEY", ""),
		BrevoBaseURL:          getEnv("BREVO_BASE_URL", "https://api.brevo.com"),
		MailFromEmail:         getEnv("MAIL_FROM_EMAIL", ""),
		MailFromName:          getEnv("MAIL_FROM_NAME", "Nox"),
		MediaUploadBaseURL:    getEnv("MEDIA_UPLOAD_BASE_URL", ""),
		MediaPublicBaseURL:    getEnv("MEDIA_PUBLIC_BASE_URL", ""),
		MediaProcessingSecret: getEnv("MEDIA_PROCESSING_SECRET", ""),
		GhostPersonaSecret:    getEnv("GHOST_PERSONA_SECRET", ""),
		Environment:           getEnv("ENVIRONMENT", getEnv("ENV", "development")),
	}

	return cfg, cfg.validate()
}

func (c *Config) validate() error {
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}
	if c.RedisURL == "" {
		return errors.New("REDIS_URL is required")
	}
	if c.JWTAccessSecret == "" {
		return errors.New("JWT_ACCESS_SECRET is required")
	}
	if c.JWTRefreshSecret == "" {
		return errors.New("JWT_REFRESH_SECRET is required")
	}
	if c.GhostPersonaSecret == "" {
		return errors.New("GHOST_PERSONA_SECRET is required")
	}
	return nil
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
