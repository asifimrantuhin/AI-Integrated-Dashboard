package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	JWTTTL        int
	APIPort       string
	AIServiceURL  string
}

var AppConfig Config

func LoadConfig() {
	// Load .env file
	_ = godotenv.Load()

	AppConfig = Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "3306"),
		DBUser:        getEnv("DB_USER", "root"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", "idash"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		JWTTTL:        getEnvInt("JWT_TTL", 24),
		APIPort:       getEnv("API_PORT", "8080"),
		AIServiceURL:  getEnv("AI_SERVICE_URL", "http://localhost:8000"),
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
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

