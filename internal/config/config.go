package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Config struct {
	TelegramToken               string
	PathToBannedWords           string
	TgNegativeChannelInviteLink string
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

	env := os.Getenv("APP_ENV")

	var path string
	if env == "production" {
		path = "internal/config/banned_words.json"
	} else {
		path = "../internal/config/banned_words.json"
	}

	negativeChatInviteLink := os.Getenv("NEGATIVE_CHAT_INVITE_LINK")

	return &Config{
		TelegramToken:               token,
		PathToBannedWords:           path,
		TgNegativeChannelInviteLink: negativeChatInviteLink,
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
