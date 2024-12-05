package bot

import (
	"log"
	"strings"

	"NegativeDetector/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

	bannedWords := []string{"упал", "реклама", "кроватки", "разбушевавшейся"}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Проверяем текст сообщения и подписи
		text := update.Message.Text
		if update.Message.Caption != "" {
			text += " " + update.Message.Caption
		}

		// Если есть запрещённые слова, удаляем сообщение
		if containsBannedWord(text, bannedWords) {
			deleteMessage(bot, update.Message)
			sendNotification(bot, update.Message)
		}
	}

	return nil
}

func containsBannedWord(text string, bannedWords []string) bool {
	text = strings.ToLower(text)
	for _, word := range bannedWords {
		if strings.Contains(text, word) {
			return true
		}
	}
	return false
}
