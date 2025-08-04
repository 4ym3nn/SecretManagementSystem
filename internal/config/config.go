package config

import (
    "os"
    "strings"
    
    "github.com/joho/godotenv"
)

type Config struct {
    Port          string
    DatabaseURL   string
    JWTSecret     string
    EncryptionKey string
    APIKeys       []string
}

func Load() *Config {
    godotenv.Load()
    
    return &Config{
        Port:          getEnv("PORT", "8080"),
        DatabaseURL:   getEnv("DATABASE_URL", "postgres://user:password@localhost/secretdb?sslmode=disable"),
        JWTSecret:     getEnv("JWT_SECRET", "your-jwt-secret-key"),
        EncryptionKey: getEnv("ENCRYPTION_KEY", "your-32-byte-encryption-key-here"),
        APIKeys:       strings.Split(getEnv("API_KEYS", "default-api-key"), ","),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
