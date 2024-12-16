package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func forwardMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	forwardMsg := tgbotapi.NewForward(negativeChatId, message.Chat.ID, message.MessageID)
	if _, err := bot.Request(forwardMsg); err != nil {
		return err
	}

	return nil
}
