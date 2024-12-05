package main

import (
	"NegativeDetector/internal/bot"
	"NegativeDetector/internal/config"
	"log"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Создание и запуск бота
	if err := bot.StartBot(cfg); err != nil {
		log.Fatalf("Error starting bot: %v", err)
	}
}
