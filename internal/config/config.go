package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  string
	DataBaseUrl string
	JwtSecret   string
}

func Load() (*Config, error) {
	godotenv.Load()
	return &Config{
		ServerPort:  os.Getenv("SERVER_PORT"),
		DataBaseUrl: os.Getenv("DATABASE_URL"),
		JwtSecret:   os.Getenv("JWT_SECRET"),
	}, nil
}
