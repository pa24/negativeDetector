package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Config struct {
	TelegramToken string
}

type WordConfig struct {
	BannedWords []string `json:"banned_words"`
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

func LoadBannedWords(filepath string) ([]string, error) {
	// Проверяем, существует ли файл

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Printf("Failed to close file %s: %v", filepath, cerr)
		}
	}()

	var config WordConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return config.BannedWords, nil
}
