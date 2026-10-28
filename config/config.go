package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort           string
	DBHost            string
	DBPort            string
	DBName            string
	DBUser            string
	DBPassword        string
	LogLevel          string
	ProductServiceURL string
}

func Load() *Config {
	return &Config{
		AppPort:           getEnv("APP_PORT", "8080"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBName:            getEnv("DB_NAME", "orderdb"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "postgres"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		ProductServiceURL: getEnv("PRODUCT_SERVICE_URL", "http://localhost:9090"),
	}
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
