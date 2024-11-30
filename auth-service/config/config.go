package config

import (
	logger "auth-service/internal"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddr  string
	DatabaseURL string
	KafkaURL    string
	RedisURL    string
}

func Load() *Config {
	err := godotenv.Load()

	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	return &Config{
		ServerAddr:  os.Getenv("SERVER_ADDR"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		KafkaURL:    os.Getenv("KAFKA_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),
	}
}
