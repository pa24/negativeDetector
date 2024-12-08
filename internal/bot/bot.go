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
	bannedWords, err = config.LoadBannedWords("../internal/config/banned_words.json")

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
	//notificationSent := make(map[string]bool)
	deletedID := make(map[int]bool)

	for update := range updates {
		exist := checkForMediaGroup(update.Message, mediaGroupCache, deletedID)
		if exist {
			deleteMediaGroup(bot, update.Message, mediaGroupCache)
			continue
		}
		if update.Message != nil {
			if update.Message.ForwardFrom != nil || update.Message.ForwardFromChat != nil {
				isNegative := wordFilter(update.Message)

				if isNegative {
					if err := deleteMessage(bot, update.Message); err != nil {
						log.Printf("Failed to delete message with id = %d in media group: %v", update.Message.MessageID, err)
					}
					deletedID[update.Message.MessageID] = true
					if update.Message.MediaGroupID != "" {
						mediaGroupCache[update.Message.MediaGroupID] = true
					}
					sendNotification(bot, update.Message)
				}
			}
		}
	}

	return nil
}

func deleteMediaGroup(bot *tgbotapi.BotAPI, message *tgbotapi.Message, mediaGroupCache map[string]bool) {
	if mediaGroupCache[message.MediaGroupID] {
		if err := deleteMessage(bot, message); err != nil {
			log.Printf("Failed to delete message with id = %d in media group: %v", message.MessageID, err) // Логируем ошибку
		}
	}
}

func checkForMediaGroup(message *tgbotapi.Message, mediaGroupCache map[string]bool, deletedID map[int]bool) bool {
	return message.MediaGroupID != "" && mediaGroupCache[message.MediaGroupID] && !deletedID[message.MessageID]
}

func prepareNotification(bot *tgbotapi.BotAPI, message *tgbotapi.Message, notificationSent map[string]bool) {

	sendNotification(bot, message)
	return
}

func wordFilter(message *tgbotapi.Message) bool {
	text := getMessageText(message)
	return containsBannedWord(text)
}

func containsBannedWord(text string) bool {
	if text == "" {
		return false
	}
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
	return update.Message.ForwardFrom != nil ||
		update.Message.ForwardFromChat != nil ||
		update.Message.Caption != "" ||
		update.Message.Text != ""
}
