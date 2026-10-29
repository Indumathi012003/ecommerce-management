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
		OrderServicePort:   os.Getenv("ORDER_SERVICE_PORT"),
		ProductServicePort: os.Getenv("PRODUCT_SERVICE_PORT"),
		OrderServiceURL:    os.Getenv("ORDER_SERVICE_URL"),
		ProductServiceURL:  os.Getenv("PRODUCT_SERVICE_URL"),
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             os.Getenv("DB_PORT"),
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
		DBSSLMode:          os.Getenv("DB_SSLMODE"),
		LogLevel:           os.Getenv("LOG_LEVEL"),
		AdminPassword:      os.Getenv("ADMIN_PASSWORD"),
		CustomerPassword:   os.Getenv("CUSTOMER_PASSWORD"),
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}
