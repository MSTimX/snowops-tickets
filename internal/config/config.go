package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config хранит все настройки сервиса.
type Config struct {
	HTTPPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBTimeZone string

	JWTSecret string

	AuthServiceURL       string
	RolesServiceURL      string
	OperationsServiceURL string
	AIServiceURL         string
}

// getEnv достаёт переменную окружения, если нет — берёт дефолт.
func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return defaultValue
}

// Load загружает .env (если есть) и собирает конфиг.
func Load() *Config {
	// Пытаемся подгрузить .env, если нет — не считаем это ошибкой.
	_ = godotenv.Load()

	cfg := &Config{
		HTTPPort: getEnv("HTTP_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "snowops_tickets"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		DBTimeZone: getEnv("DB_TIMEZONE", "Asia/Almaty"),

		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-me"),

		AuthServiceURL:       getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		RolesServiceURL:      getEnv("ROLES_SERVICE_URL", "http://localhost:8082"),
		OperationsServiceURL: getEnv("OPERATIONS_SERVICE_URL", "http://localhost:8083"),
		AIServiceURL:         getEnv("AI_SERVICE_URL", ""),
	}

	log.Printf("config loaded: HTTP_PORT=%s, DB_HOST=%s, DB_NAME=%s", cfg.HTTPPort, cfg.DBHost, cfg.DBName)

	return cfg
}
