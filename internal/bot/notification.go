package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

//const textNotification = "Сообщение удалено: обнаружен негатив, посмотреть его можно тут: "

// sendNotification отправляет уведомление о том, что сообщение было удалено
func sendNotification(bot *tgbotapi.BotAPI, message *tgbotapi.Message, mediaGroupCache map[string]bool, link string) {
	if message.MediaGroupID != "" {
		mediaGroupCache[message.MediaGroupID] = true
	}
	messageFromUserName := message.From.UserName
	textNotification := fmt.Sprintf("Сообщение от %s удалено: обнаружен негатив. Посмотреть его можно тут: %s", messageFromUserName, link)
	msg := tgbotapi.NewMessage(message.Chat.ID, textNotification)
	//msg.ParseMode = "Markdown"
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send notification: %v", err)
	}
}
