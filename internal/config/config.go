package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBHost         string
	DBPort         string
	DBName         string
	DBUser         string
	DBPass         string
	JWTSecret      string
	RedisURL       string
	MigrateOnStart bool
}

// Load reads configuration from environment variables and falls back to
// sensible defaults when values are missing. This prevents runtime failures
// when env vars are not set during local development.
func Load() Config {
	// try to load .env (if present) so os.Getenv can read values from it
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "tamagoam"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}

	dbPass := os.Getenv("DB_PASS")
	// dbPass may be empty; keep as-is

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "tamagoam-dev-secret-change-in-production"
	}

	migrate := false
	m := os.Getenv("MIGRATE_ON_START")
	if m != "" {
		m = strings.ToLower(strings.TrimSpace(m))
		if m == "true" || m == "1" || m == "yes" {
			migrate = true
		}
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://127.0.0.1:6379/0"
	}

	return Config{
		Port:           port,
		DBHost:         dbHost,
		DBPort:         dbPort,
		DBName:         dbName,
		DBUser:         dbUser,
		DBPass:         dbPass,
		JWTSecret:      jwtSecret,
		RedisURL:       redisURL,
		MigrateOnStart: migrate,
	}
}
