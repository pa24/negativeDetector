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
	bannedWords, err = config.LoadBannedWords("internal/config/banned_words.json")
	if err != nil {
		return err
	}
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return err
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	mediaGroupCache := make(map[string]bool)

	for update := range updates {
		if mediaGroupCache[update.Message.MediaGroupID] {
			deleteMediaGroup(bot, update.Message)
			continue
		}
		if !isMessageValid(&update) {
			continue
		}
		if wordFilter(update.Message) {
			handleNegativeMessage(bot, update.Message, mediaGroupCache)
		}
	}

	return nil
}

func handleNegativeMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, mediaGroupCache map[string]bool) {
	if err := deleteMessage(bot, message); err != nil {
		log.Printf("Failed to delete message with id = %d in media group: %v", message.MessageID, err)
		return
	}
	if message.MediaGroupID != "" {
		mediaGroupCache[message.MediaGroupID] = true
	}

	sendNotification(bot, message)
}

func deleteMediaGroup(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if err := deleteMessage(bot, message); err != nil {
		log.Printf("Failed to delete message with id = %d in media group: %v", message.MessageID, err) // Логируем ошибку
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
