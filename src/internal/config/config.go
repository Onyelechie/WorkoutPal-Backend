package config

import (
	"os"
)

type Config struct {
	DatabaseURL    string
	GoogleClientID string
	JWTSecret      string
	Port           string
}

func Load() *Config {
	return &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "host=127.0.0.1 port=5432 user=user password=password dbname=workoutpal sslmode=disable"),
		GoogleClientID: getEnv("GOOGLE_CLIENT_ID", ""),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
		Port:           getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}