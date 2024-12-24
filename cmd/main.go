package main

import (
	negativeDetector "NegativeDetector"
	"NegativeDetector/internal/config"
	"NegativeDetector/internal/database"
	"NegativeDetector/internal/database/migrations"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	// Подключаемся к базе данных
	db, err := database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	//Запуск миграций
	if err := migrations.RunMigrations(db.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations applied successfully")

	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("can't creates a new BotAPI instance: %w", err)
	}
	botAPI.Debug = false
	log.Printf("Authorized on account %s", botAPI.Self.UserName)
	log.WithFields(log.Fields{
		"username": botAPI.Self.UserName,
	}).Info("Bot successfully authorized")

	log.Println("запуск бота прошел успешно")
	//создание сервера
	go negativeDetector.StartServer(botAPI, db)
}
