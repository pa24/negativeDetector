package main

import (
	"NegativeDetector/internal/bot"
	"NegativeDetector/internal/config"
	"log"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Создание и запуск бота
	if err = bot.StartBot(cfg); err != nil {
		log.Fatalf("Error starting bot: %v", err)
	}
}
