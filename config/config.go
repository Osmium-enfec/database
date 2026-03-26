package config

import (
	"net/url"
	"os"
	"strings"
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
	URL      string // For Render's DATABASE_URL
}

type RedisConfig struct {
	Host     string `json:"host" ini:"host"`
	Port     string `json:"port" ini:"port"`
	User     string `json:"user" ini:"user"`
	Password string `json:"password" ini:"password"`
	Database int    `json:"database" ini:"database"`
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

type CORSConfig struct {
	AllowedOrigins []string
}

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
	CORS     CORSConfig
}

// ParseDatabaseURL parses Render's DATABASE_URL format: postgres://user:password@host:port/dbname
func parseDatabaseURL(dbURL string) (host, port, user, password, dbname, sslmode string) {
	if dbURL == "" {
		return "localhost", "5432", "postgres", "postgres", "content_review", "disable"
	}

	// Handle both postgres:// and postgresql:// URLs
	if strings.HasPrefix(dbURL, "postgresql://") {
		dbURL = "postgres://" + dbURL[13:]
	}

	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		// Fallback to environment variables if parsing fails
		return "localhost", "5432", "postgres", "postgres", "content_review", "disable"
	}

	user = parsedURL.User.Username()
	password, _ = parsedURL.User.Password()
	host = parsedURL.Hostname()
	port = parsedURL.Port()
	if port == "" {
		port = "5432"
	}

	// Extract database name from path (e.g., "/dbname" -> "dbname")
	dbname = strings.TrimPrefix(parsedURL.Path, "/")
	if dbname == "" {
		dbname = "content_review"
	}

	// Check for sslmode in query parameters
	sslmode = parsedURL.Query().Get("sslmode")
	if sslmode == "" {
		sslmode = "require" // Render typically requires SSL
	}

	return
}

func LoadConfig() *Config {
	// Check for Render's DATABASE_URL first
	dbURL := os.Getenv("DATABASE_URL")
	var dbConfig DatabaseConfig

	if dbURL != "" {
		// Parse Render's DATABASE_URL
		host, port, user, password, dbname, sslmode := parseDatabaseURL(dbURL)
		dbConfig = DatabaseConfig{
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
			DBName:   dbname,
			SSLMode:  sslmode,
			MaxConn:  getEnvInt("DB_MAX_CONN", 25),
			URL:      dbURL,
		}
	} else {
		// Fall back to individual environment variables (for local development)
		dbConfig = DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "content_review"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			MaxConn:  getEnvInt("DB_MAX_CONN", 25),
		}
	}

	return &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "content-review-api"),
			Environment: getEnv("ENVIRONMENT", "production"),
			Debug:       getEnv("DEBUG", "false") == "true",
		},
		Database: dbConfig,
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-change-in-production"),
			ExpiryHours:   getEnvInt("JWT_EXPIRY_HOURS", 720), // 30 days
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		CORS: CORSConfig{
			AllowedOrigins: parseCORSOrigins(
				getEnv("CORS_ALLOWED_ORIGINS",
					"https://content-review-api-bnkf.onrender.com,https://id-preview--e5b904ce-9f96-4c37-9e1a-41a95d44462a.lovable.app,http://localhost:3000,http://localhost:5173,https://localhost:3000"),
			),
		},
	}
}

// parseCORSOrigins parses comma-separated CORS origins
func parseCORSOrigins(originsStr string) []string {
	if originsStr == "" {
		return []string{
			// Lovable URLs
			"https://e5b904ce-9f96-4c37-9e1a-41a95d44462a.lovableproject.com",
			"https://id-preview--e5b904ce-9f96-4c37-9e1a-41a95d44462a.lovable.app",
			"https://contentenfeca.lovable.app",
			// Local development
			"http://localhost:3000",
			"http://localhost:5173",
			"https://localhost:3000",
		}
	}

	origins := strings.Split(originsStr, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return origins
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
