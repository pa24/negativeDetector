package main

import (
	"NegativeDetector/internal/bot"
	"NegativeDetector/internal/config"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	// Загрузка конфигурации
	log.SetFormatter(&log.JSONFormatter{}) // Логи в JSON формате
	log.SetLevel(log.InfoLevel)            // Уровень логирования по умолчанию
	log.SetOutput(os.Stdout)               // Выводим логи в консоль

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Создание и запуск бота
	if err = bot.StartBot(cfg); err != nil {
		log.Fatalf("Error starting bot: %v", err)
	}
}
