package config

import (
	"fmt"
	"os"
	"strconv"

	"effective-mobile/internal/models"

	"github.com/joho/godotenv"
)

func Load() (*models.Config, error) {

	if err := godotenv.Load("ini.env"); err != nil {
		return nil, fmt.Errorf("error loading ini.env: %w", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("environment variable PORT is not set")
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		return nil, fmt.Errorf("environment variable DB_HOST is not set")
	}

	portStr := os.Getenv("DB_PORT")
	if portStr == "" {
		return nil, fmt.Errorf("environment variable DB_PORT is not set")
	}
	dbPort, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT %q: %w", portStr, err)
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		return nil, fmt.Errorf("environment variable DB_USER is not set")
	}
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		return nil, fmt.Errorf("environment variable DB_PASSWORD is not set")
	}
	name := os.Getenv("DB_NAME")
	if name == "" {
		return nil, fmt.Errorf("environment variable DB_NAME is not set")
	}

	cfg := &models.Config{
		Port: port,
		DBConfig: models.DBConfig{
			Host:     host,
			Port:     dbPort,
			User:     user,
			Password: pass,
			DbName:   name,
		},
	}

	return cfg, nil
}
