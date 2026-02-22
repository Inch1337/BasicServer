package config

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	ServerPort string
}

var ErrInvalidConfig = errors.New("invalid config")

func New() (*Config, error) {
	cfg := &Config{
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "postgres"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		ServerPort: getEnv("SERVER_PORT", ":8081"),
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Validate() error {
	if c.DBHost == "" {
		return errors.Join(ErrInvalidConfig, errors.New("DB_HOST is required"))
	}
	if c.DBPort == "" {
		return errors.Join(ErrInvalidConfig, errors.New("DB_PORT is required"))
	}
	if _, err := strconv.Atoi(c.DBPort); err != nil {
		return errors.Join(ErrInvalidConfig, errors.New("DB_PORT must be a number"))
	}
	if c.ServerPort != "" && !strings.HasPrefix(c.ServerPort, ":") {
		c.ServerPort = ":" + c.ServerPort
	}
	return nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
