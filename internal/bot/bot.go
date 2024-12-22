package bot

import (
	"NegativeDetector/internal/config"
	"NegativeDetector/internal/database"
	"NegativeDetector/internal/database/migrations"
	"NegativeDetector/internal/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"strings"
)

var bannedWords []string

// StartBot инициализирует и запускает бота
func StartBot(cfg *config.Config) error {
	var err error

	// Подключаемся к базе данных
	db, err := database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Запуск миграций
	if err := migrations.RunMigrations(db.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations applied successfully")

	bannedWords, err = config.LoadBannedWords(cfg.PathToBannedWords)
	if err != nil {
		return err
	}
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return err
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)
	log.WithFields(log.Fields{
		"username": bot.Self.UserName,
	}).Info("Bot successfully authorized")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	mediaGroupCache := make(map[string]bool)

	for update := range updates {
		if update.Message != nil {
			handlers.SaveMessageHandler(db, update.Message)
		}

		if isMessageGroup(update.Message, mediaGroupCache, cfg.TargetChatID) {
			forwardAndDelete(bot, update.Message, cfg.ForwardChatID)
			continue
		}

		if !isMessageValid(&update) {
			continue
		}
		if wordFilter(update.Message) {
			log.WithFields(log.Fields{
				"user_id":    update.Message.From.ID,
				"user_name":  update.Message.From.UserName,
				"message_id": update.Message.MessageID,
				"chat_id":    update.Message.Chat.ID,
			}).Warn("Message contains banned words")

			forwardAndDelete(bot, update.Message, cfg.ForwardChatID)
			sendNotification(bot, update.Message, mediaGroupCache, cfg.TgNegativeChannelInviteLink)
		}
	}

	return nil
}

func isMessageGroup(message *tgbotapi.Message, mediaGroupCache map[string]bool, targetChatID int64) bool {
	if message == nil { // Проверяем, что Message не nil
		return false
	}
	if message.MediaGroupID != "" && message.Chat.ID == targetChatID {
		if mediaGroupCache[message.MediaGroupID] {
			log.WithFields(log.Fields{
				"media_group_id": message.MediaGroupID,
				"chat_id":        message.Chat.ID,
			}).Debug("Message is part of an existing media group")
			return true
		}
	}
	return false
}

func forwardAndDelete(bot *tgbotapi.BotAPI, message *tgbotapi.Message, forwardChatID int64) {
	if err := forwardMessage(bot, message, forwardChatID); err != nil {
		log.WithFields(log.Fields{
			"message_id": message.MessageID,
			"chat_id":    message.Chat.ID,
			"error":      err,
		}).Error("Failed to forward message")

	}
	if err := deleteMessage(bot, message); err != nil {
		log.WithFields(log.Fields{
			"message_id": message.MessageID,
			"chat_id":    message.Chat.ID,
			"error":      err,
		}).Error("Failed to delete message")
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
	if update.Message == nil {
		return false
	}
	valid := (update.Message.Caption != "" || update.Message.Text != "") &&
		(update.Message.ForwardFrom != nil || update.Message.ForwardFromChat != nil)
	if !valid {
		log.WithFields(log.Fields{
			"update": update,
		}).Debug("Message is invalid or empty")
	}
	return valid
}
