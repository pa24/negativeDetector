package bot

import (
	"log"
	"strings"

	"NegativeDetector/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bannedWords = []string{"упал", "реклама", "кроватки", "разбушевавшейся", "готовить"}

// StartBot инициализирует и запускает бота
func StartBot(cfg *config.Config) error {
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
	notificationSent := make(map[string]bool)
	deletedID := make(map[int]bool)

	for update := range updates {
		if update.Message != nil {
			if update.Message.ForwardFrom != nil || update.Message.ForwardFromChat != nil {
				isNegative := wordFilter(update.Message)

				if isNegative {
					deleteMessage(bot, update.Message)
					deletedID[update.Message.MessageID] = true
					if update.Message.MediaGroupID != "" {
						mediaGroupCache[update.Message.MediaGroupID] = true
					}
				}
				exist := checkForMediaGroup(update.Message, mediaGroupCache, deletedID)
				if exist {
					deleteMediaGroup(bot, update.Message, mediaGroupCache)

				}
				if isNegative || exist {
					prepareNotification(bot, update.Message, notificationSent)
				}
			}
		}
	}

	return nil
}

func deleteMediaGroup(bot *tgbotapi.BotAPI, message *tgbotapi.Message, mediaGroupCache map[string]bool) {
	if mediaGroupCache[message.MediaGroupID] {
		deleteMessage(bot, message)
	}
}

func checkForMediaGroup(message *tgbotapi.Message, mediaGroupCache map[string]bool, deletedID map[int]bool) bool {
	if message.MediaGroupID != "" && mediaGroupCache[message.MediaGroupID] {
		if !deletedID[message.MessageID] {
			return true
		}
	}
	return false
}

func prepareNotification(bot *tgbotapi.BotAPI, message *tgbotapi.Message, notificationSent map[string]bool) {
	if message.MediaGroupID != "" {
		if !notificationSent[message.MediaGroupID] {
			sendNotification(bot, message)
			notificationSent[message.MediaGroupID] = true
			return
		}
	} else {
		sendNotification(bot, message)
		return
	}
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
