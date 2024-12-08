package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// deleteMessage удаляет текстовые и медиа-сообщения
func deleteMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	deleteMsg := tgbotapi.DeleteMessageConfig{
		ChatID:    message.Chat.ID,
		MessageID: message.MessageID,
	}
	if _, err := bot.Request(deleteMsg); err != nil {
		return err
	}
	return nil
}
