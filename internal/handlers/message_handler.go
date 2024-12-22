package handlers

import (
	"NegativeDetector/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5" // Замените на вашу библиотеку для работы с Telegram API
	log "github.com/sirupsen/logrus"
	"strings"
)

// SaveMessageHandler обрабатывает и сохраняет сообщение в базу данных
func SaveMessageHandler(db *database.Database, message *tgbotapi.Message) {
	// Определяем тип контента
	content := getMessageText(message)
	var contentType string
	switch {
	case message.Photo != nil:
		contentType = "photo"
	case message.Voice != nil:
		contentType = "voice"
	case message.Video != nil:
		contentType = "video"
	case message.VideoNote != nil:
		contentType = "video_note"
	default:
		contentType = "text"
	}

	// Сохраняем сообщение в базу данных
	err := db.SaveMessage(
		message.From.ID,
		message.From.UserName,
		contentType,
		content,
		message.Chat.ID,
	)
	if err != nil {
		log.Printf("Failed to save message: %v", err)
	}
}

// getMessageText извлекает текст сообщения
func getMessageText(message *tgbotapi.Message) string {
	if message.Text != "" {
		return strings.TrimSpace(message.Text)
	}
	return "<empty>"
}
