package main

import (
	"NegativeDetector/internal/bot"
	"NegativeDetector/internal/config"
	"NegativeDetector/internal/database"
	"NegativeDetector/internal/database/migrations"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
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

	// Запуск миграций
	if err := migrations.RunMigrations(db.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations applied successfully")

	newBotAPI, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("can't creates a new BotAPI instance: %w", err)
	}
	newBotAPI.Debug = false
	log.Printf("Authorized on account %s", newBotAPI.Self.UserName)
	log.WithFields(log.Fields{
		"username": newBotAPI.Self.UserName,
	}).Info("Bot successfully authorized")

	log.Println("запуск бота прошел успешно")
	//создание сервера
	go func() {
		http.HandleFunc("/send-daily-stats", func(w http.ResponseWriter, r *http.Request) {
			log.Println("сервер запускается")
			chatIDStr := r.URL.Query().Get("chat_id")
			if chatIDStr == "" {
				log.Println("chat_id is missing")
				http.Error(w, "chat_id is required", http.StatusBadRequest)
				return
			}
			chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
			if err != nil {
				log.Println("chat_id is missing: ", chatID)
				http.Error(w, "Invalid chat_id", http.StatusBadRequest)
				return
			}

			err = bot.SendDailyStats(newBotAPI, db, chatID)
			if err != nil {
				log.Errorf("Error sending daily stats: %v", err)
				http.Error(w, "Failed to send stats", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Stats sent successfully"))
		})

		log.Info("Starting server on :8080")
		err = http.ListenAndServe(":8080", nil)
		if err != nil {
			return
		}
	}()

	//Создание и запуск бота
	if err = bot.StartBot(cfg, db, newBotAPI); err != nil {
		log.Fatalf("Error starting newBotAPI: %v", err)
	}
}
