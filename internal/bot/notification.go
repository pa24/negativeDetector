package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const textNotification = "Сообщение удалено: обнаружен негатив."

// sendNotification отправляет уведомление о том, что сообщение было удалено
func sendNotification(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, textNotification)
	msg.ParseMode = "Markdown"
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send notification: %v", err)
	}
}
