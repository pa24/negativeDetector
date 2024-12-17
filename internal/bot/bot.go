package bot

import (
	"log"
	"strings"

	"NegativeDetector/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bannedWords []string

// StartBot инициализирует и запускает бота
func StartBot(cfg *config.Config) error {
	var err error
	bannedWords, err = config.LoadBannedWords(cfg.PathToBannedWords)
	if err != nil {
		return err
	}
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return err
	}

	targetChatID := int64(cfg.TargetChatID)
	forwardChatID := int64(cfg.ForwardChatID)

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	mediaGroupCache := make(map[string]bool)

	for update := range updates {
		if isMessageGroup(update.Message, mediaGroupCache, targetChatID) {
			forwardMediaGroup(bot, update.Message, forwardChatID)
			continue
		}

		if !isMessageValid(&update) {
			continue
		}
		if wordFilter(update.Message) {
			handleNegativeMessage(bot, update.Message, forwardChatID)
			sendNotification(bot, update.Message, mediaGroupCache, cfg.TgNegativeChannelInviteLink)
		}
	}

	return nil
}

func isMessageGroup(message *tgbotapi.Message, mediaGroupCache map[string]bool, targetChatID int64) bool {
	if message != nil && message.MediaGroupID != "" && message.Chat.ID == targetChatID {
		if mediaGroupCache[message.MediaGroupID] {
			return true
		}
	}
	return false
}

func handleNegativeMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, forwardChatID int64) {
	if err := forwardMessage(bot, message, forwardChatID); err != nil {
		log.Printf("Failed to forward message with id = %d in media group: %v", message.MessageID, err)
		return
	}
	if err := deleteMessage(bot, message); err != nil {
		log.Printf("Failed to delete message with id = %d in media group: %v", message.MessageID, err) // Логируем ошибку
	}
}

func forwardMediaGroup(bot *tgbotapi.BotAPI, message *tgbotapi.Message, forwardChatID int64) {
	if err := forwardMessage(bot, message, forwardChatID); err != nil {
		log.Printf("Failed to forward message with id = %d in media group: %v", message.MessageID, err)
	}
	if err := deleteMessage(bot, message); err != nil {
		log.Printf("Failed to delete message with id = %d in media group: %v", message.MessageID, err)
	}

}

func wordFilter(message *tgbotapi.Message) bool {
	text := getMessageText(message)
	return containsBannedWord(text)
}

func containsBannedWord(text string) bool {
	text = strings.ToLower(text)
	for _, word := range bannedWords {
		if strings.Contains(text, word) {
			return true
		}
	}
	return false
}

func getMessageText(message *tgbotapi.Message) string {
	if message.Caption != "" {
		return message.Caption
	}
	return message.Text
}

func isMessageValid(update *tgbotapi.Update) bool {
	return (update.Message.Caption != "" || update.Message.Text != "") &&
		(update.Message.ForwardFrom != nil || update.Message.ForwardFromChat != nil)
}
