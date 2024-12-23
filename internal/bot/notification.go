package bot

import (
	"NegativeDetector/internal/database"
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

func SendDailyStats(bot *tgbotapi.BotAPI, db *database.Database, chatID int64) error {
	stats, err := database.GetDailyStats(db, chatID)
	if err != nil {
		return fmt.Errorf("failed to get daily stats: %v", err)
	}

	// Форматируем сообщение
	message := fmt.Sprintf(
		"📊 *Ежедневная статистика чата за вчерашний день:*\n\n"+
			"💬 Всего сообщений было отправлено: *%d*\n"+
			"📝 Всего слов было напечатано: *%d*\n"+
			"👑 Самый активный пользователь: *%s* (сообщений: %d)\n"+
			"🎙️ Голосовых сообщений больше всего отправил: *%s* (%d)\n"+
			"🎥 Видеосообщений больше всего отправил: *%s* (%d)\n",
		stats.TotalMessages,
		stats.TotalWords,
		stats.MostActiveUser, stats.MostActiveUserMessages,
		stats.TopVoiceUser, stats.TopVoiceMessages,
		stats.TopVideoUser, stats.TopVideoMessages,
	)

	// Создаем сообщение для Telegram
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"

	// Отправляем сообщение
	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}
