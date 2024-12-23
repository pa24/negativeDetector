package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
)

type Config struct {
	TelegramToken               string
	TgNegativeChannelInviteLink string
	TargetChatID                int64
	ForwardChatID               int64
	DatabaseURL                 string
	Enviroment                  string
}

type wordConfig struct {
	BannedWords []string `json:"banned_words"`
}

// LoadConfig загружает конфигурацию из переменных окружения или других источников
func LoadConfig() (*Config, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, errors.New("TELEGRAM_BOT_TOKEN is not set")
	}
	env := os.Getenv("APP_ENV")

	negativeChatInviteLink := os.Getenv("NEGATIVE_CHAT_INVITE_LINK")
	targetChatIDStr := os.Getenv("TARGET_CHAT_ID")
	forwardChatIDStr := os.Getenv("FORWARD_CHAT_ID")

	dataBaseURL := os.Getenv("DATABASE_URL")

	targetChatID, err := strconv.Atoi(targetChatIDStr)
	if err != nil {
		log.Printf("Failed to convert TARGET_CHAT_ID to int: %v", err)
		return nil, err
	}

	forwardChatID, err := strconv.Atoi(forwardChatIDStr)
	if err != nil {
		log.Printf("Failed to convert FORWARD_CHAT_ID to int: %v", err)
		return nil, err
	}

	return &Config{
		TelegramToken:               token,
		TgNegativeChannelInviteLink: negativeChatInviteLink,
		TargetChatID:                int64(targetChatID),
		ForwardChatID:               int64(forwardChatID),
		DatabaseURL:                 dataBaseURL,
		Enviroment:                  env,
	}, nil
}

func LoadBannedWords(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Printf("Failed to close file %s: %v", filepath, cerr)
		}
	}()

	var config wordConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return config.BannedWords, nil
}
