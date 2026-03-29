package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	GinMode           string
	JWTSecret         string
	JWTExpiryHours    int
	OpenAIKey         string
	OpenAIModel       string
	DBPath            string
	RateLimitRequests int
	RateLimitPeriod   string
}

var App *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	rateLimit, _ := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "100"))

	App = &Config{
		Port:              getEnv("PORT", "8080"),
		GinMode:           getEnv("GIN_MODE", "debug"),
		JWTSecret:         getEnv("JWT_SECRET", "fallback-secret-change-me"),
		JWTExpiryHours:    jwtExpiry,
		OpenAIKey:         getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:       getEnv("OPENAI_MODEL", "gpt-3.5-turbo"),
		DBPath:            getEnv("DB_PATH", "./smarttask.db"),
		RateLimitRequests: rateLimit,
		RateLimitPeriod:   getEnv("RATE_LIMIT_PERIOD", "1m"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}