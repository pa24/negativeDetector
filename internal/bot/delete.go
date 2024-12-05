package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// deleteMessage удаляет текстовые и медиа-сообщения
func deleteMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	deleteMsg := tgbotapi.DeleteMessageConfig{
		ChatID:    message.Chat.ID,
		MessageID: message.MessageID,
	}
	if _, err := bot.Request(deleteMsg); err != nil {
		log.Printf("Failed to delete message: %v", err)
	}

	// Удаление связанных медиа-сообщений
	if message.Photo != nil || message.Video != nil || message.Document != nil {
		mediaMessageID := message.MessageID + 1
		deleteMediaMsg := tgbotapi.DeleteMessageConfig{
			ChatID:    message.Chat.ID,
			MessageID: mediaMessageID,
		}
		if _, err := bot.Request(deleteMediaMsg); err != nil {
			log.Printf("Failed to delete media message: %v", err)
		}
	}
}
