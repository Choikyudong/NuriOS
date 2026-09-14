package config

import (
	"os"
)

type Config struct {
	DBUrl     string
	JWTSecret string
	Port      string
}

func Load() Config {
	return Config{
		DBUrl:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/nurios?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-me"),
		Port:      getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
