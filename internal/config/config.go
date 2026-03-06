package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        int
	DatabaseURL string
	Environment string
	OpenAIKey   string
}

func LoadConfig() Config {
	port := 8080
	if v, ok := os.LookupEnv("PORT"); ok {
		if parsed, err := strconv.Atoi(v); err == nil {
			port = parsed
		}
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/retrotriage?sslmode=disable"
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	return Config{
		Port:        port,
		DatabaseURL: databaseURL,
		Environment: environment,
		OpenAIKey:   os.Getenv("OPENAI_API_KEY"),
	}
}
