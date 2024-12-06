package config

import (
	"errors"
	"os"
)

type Config struct {
	TelegramToken string
}

// LoadConfig загружает конфигурацию из переменных окружения или других источников
func LoadConfig() (*Config, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, errors.New("TELEGRAM_BOT_TOKEN is not set")
	}

	return &Config{
		TelegramToken: token,
	}, nil
}
