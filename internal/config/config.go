package config

import (
	"log"
	"os"
)

type Config struct {
	TelegramToken string
}

// LoadConfig загружает конфигурацию из переменных окружения или других источников
func LoadConfig() *Config {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	return &Config{
		TelegramToken: token,
	}
}
