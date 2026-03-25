package config

import (
	"os"
)

type AppConfig struct {
	Name        string
	Environment string
	Debug       bool
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConn  int
}

type JWTConfig struct {
	Secret        string
	RefreshSecret string
	ExpiryHours   int
}

type ServerConfig struct {
	Host string
	Port string
}

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
}

func LoadConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "content-review-api"),
			Environment: getEnv("ENVIRONMENT", "development"),
			Debug:       getEnv("DEBUG", "false") == "true",
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "content_review"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			MaxConn:  getEnvInt("DB_MAX_CONN", 25),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-change-in-production"),
			ExpiryHours:   getEnvInt("JWT_EXPIRY_HOURS", 720), // 30 days
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		// Parse the integer
		intVal := 0
		for _, c := range value {
			if c >= '0' && c <= '9' {
				intVal = intVal*10 + int(c-'0')
			} else {
				return defaultValue
			}
		}
		return intVal
	}
	return defaultValue
}
