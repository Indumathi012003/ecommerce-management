package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GatewayPort        string
	OrderServicePort   string
	ProductServicePort string
	OrderServiceURL    string
	ProductServiceURL  string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	LogLevel           string
	AdminPassword      string
	CustomerPassword   string
}

func envOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func Load() Config {
	_ = godotenv.Load()

	gwPort := os.Getenv("APP_PORT")
	if gwPort == "" {
		gwPort = envOrDefault("GATEWAY_PORT", "8080")
	}

	return Config{
		GatewayPort:        gwPort,
		OrderServicePort:   envOrDefault("ORDER_SERVICE_PORT", "8081"),
		ProductServicePort: envOrDefault("PRODUCT_SERVICE_PORT", "8082"),
		OrderServiceURL:    envOrDefault("ORDER_SERVICE_URL", "http://localhost:8081"),
		ProductServiceURL:  envOrDefault("PRODUCT_SERVICE_URL", "http://localhost:8082"),
		DBHost:             envOrDefault("DB_HOST", "localhost"),
		DBPort:             envOrDefault("DB_PORT", "5432"),
		DBUser:             envOrDefault("DB_USER", "postgres"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             envOrDefault("DB_NAME", "orderdb"),
		DBSSLMode:          envOrDefault("DB_SSLMODE", "disable"),
		LogLevel:           envOrDefault("LOG_LEVEL", "info"),
		AdminPassword:      envOrDefault("ADMIN_PASSWORD", "admin123"),
		CustomerPassword:   envOrDefault("CUSTOMER_PASSWORD", "customer123"),
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}
