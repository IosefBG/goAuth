package config

import (
	"os"
)

// AppConfig holds the configuration settings for the application.
type AppConfig struct {
	Environment string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
}

// LoadConfig loads the application configuration from environment variables.
func LoadConfig() AppConfig {
	return AppConfig{
		Environment: os.Getenv("ENVIRONMENT"),
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
	}
}
