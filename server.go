package negativeDetector

import (
	"NegativeDetector/internal/bot"
	"NegativeDetector/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// StartServer запускает HTTP-сервер для обработки запросов.
func StartServer(botAPI *tgbotapi.BotAPI, db *database.Database) {
	http.HandleFunc("/send-daily-stats", func(w http.ResponseWriter, r *http.Request) {
		chatIDStr := r.URL.Query().Get("chat_id")
		if chatIDStr == "" {
			log.Println("chat_id is missing")
			http.Error(w, "chat_id is required", http.StatusBadRequest)
			return
		}

		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			log.Printf("Invalid chat_id: %v", err)
			http.Error(w, "Invalid chat_id", http.StatusBadRequest)
			return
		}

		err = bot.SendDailyStats(botAPI, db, chatID)
		if err != nil {
			log.Errorf("Error sending daily stats: %v", err)
			http.Error(w, "Failed to send stats", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Stats sent successfully"))
	})

	log.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
