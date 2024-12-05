package bot

import (
	"log"
	"strings"

	"NegativeDetector/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bannedWords = []string{"упал", "реклама", "кроватки", "разбушевавшейся"}

// StartBot инициализирует и запускает бота
func StartBot(cfg *config.Config) error {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return err
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	mediaGroupCache := make(map[string]bool)
	notificationSent := make(map[string]bool)

	for update := range updates {
		if update.Message != nil {
			if update.Message.ForwardFrom != nil || update.Message.ForwardFromChat != nil {
				prepareToDelete(bot, update.Message, mediaGroupCache)
				prepareNotification(bot, update.Message, notificationSent)

			}
		}
	}

	return nil
}

func prepareNotification(bot *tgbotapi.BotAPI, message *tgbotapi.Message, notificationSent map[string]bool) {
	if !notificationSent[message.MediaGroupID] {
		sendNotification(bot, message)
		notificationSent[message.MediaGroupID] = true
	}
}

func prepareToDelete(bot *tgbotapi.BotAPI, message *tgbotapi.Message, cache map[string]bool) {
	mediaGroupID := message.MediaGroupID
	// Если есть запрещённые слова, удаляем сообщение
	var containsWord bool
	if mediaGroupID != "" {
		// Проверяем или обновляем статус группы
		//todo отделить проверку слов от проверки кеша. Проверка кеша должна быть внутри
		containsWord = containsBannedWord(message.Caption) || cache[mediaGroupID]
		cache[mediaGroupID] = containsWord
	}
	if containsWord {
		deleteMessage(bot, message)
	}
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
