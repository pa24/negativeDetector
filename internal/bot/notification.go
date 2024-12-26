package bot

import (
	"NegativeDetector/internal/database"
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"text/template"
)

var textUserPhrases = []string{
	"Самый активный пользователь",
	"Самый болтливый",
	"Неугомон чата",
	"Король чата",
	"Чат-лидер",
	"Главный словесный ниндзя",
	"Мастер клавиатуры",
	"Король клавиш",
	"Властелин сообщений",
	"Гуру общения",
	"Победитель в номинации 'Болтазавр-яичница'",
	"Чат-магнат",
	"Батя чатика",
	"Сердце беседы",
	"Рыцарь словесной битвы",
	"Заводила чатика",
	"Душа компании",
}

var voiceUserPhrases = []string{
	"Мастер подкастов",
	"Голос чата",
	"Король эфира",
	"Главный диктор",
	"Гуру голосовых сообщений",
	"Властелин микрофона",
	"Главный шептатель в микрофон",
	"Маэстро микрофона",
	"Голосовой гигант",
	"Гуру голосогого сообщения",
}

var videoNoneUserPhrases = []string{
	"Кинорежиссёр чата",
	"Звезда видеосообщений",
	"Король-видеомагнат",
	"Главный оператор",
	"Мастер клипов",
	"Главный кружок мейкер",
	"Старший кружкоделатель",
	"Камера-маньяк",
}

type Stats struct {
	TotalMessages             int
	TotalWords                int
	MostActiveUserPhrase      string
	MostActiveVoicePhrase     string
	MostActiveVideoNotePhrase string
	MostActiveUserID          string
	MostActiveUserMessages    string
	TopVoiceUser              string
	TopVoiceMessages          int
	TopVideoUser              string
	TopVideoMessages          int
}

// sendNotification отправляет уведомление о том, что сообщение было удалено

func SendDailyStats(bot *tgbotapi.BotAPI, db *database.Database, chatID int64, notifierChatID int64) error {
	stats, err := database.GetDailyStats(db, chatID)
	if err != nil {
		return fmt.Errorf("failed to get daily stats: %v", err)
	}

	prepareStats := Stats{
		TotalMessages:             stats.TotalMessages,
		TotalWords:                stats.TotalWords,
		MostActiveUserPhrase:      getRandomPhrase(textUserPhrases),
		MostActiveVoicePhrase:     getRandomPhrase(voiceUserPhrases),
		MostActiveVideoNotePhrase: getRandomPhrase(videoNoneUserPhrases),
		MostActiveUserID:          stats.MostActiveUserID,
		MostActiveUserMessages:    pluralizeMessages(stats.MostActiveUserMessages),
		TopVoiceUser:              stats.TopVoiceUser,
		TopVoiceMessages:          stats.TopVoiceMessages,
		TopVideoUser:              stats.TopVideoUser,
		TopVideoMessages:          stats.TopVideoMessages,
	}

	message, err := generateMessage(prepareStats)
	if err != nil {
		fmt.Println("Error generating message:", err)
		return err
	}
	// Создаем сообщение для Telegram
	msg := tgbotapi.NewMessage(notifierChatID, message)
	msg.ParseMode = "Markdown"

	// Отправляем сообщение
	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}

func generateMessage(stats Stats) (string, error) {
	templateText := `📊 *Ежедневная статистика чата за вчерашний день:*

💬 Всего сообщений было отправлено: *{{.TotalMessages}}*
📝 Всего слов было напечатано: *{{.TotalWords}}*
👑 {{.MostActiveUserPhrase}}: *{{.MostActiveUserID}}* ({{.MostActiveUserMessages}})
🎙️ {{.MostActiveVoicePhrase}}: *{{.TopVoiceUser}}* ({{.TopVoiceMessages}})
🎥 {{.MostActiveVideoNotePhrase}}: *{{.TopVideoUser}}* ({{.TopVideoMessages}})
`

	tmpl, err := template.New("stats").Parse(templateText)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, stats)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func pluralizeMessages(num int) string {
	if num%10 == 1 && num%100 != 11 {
		return fmt.Sprintf("%d сообщение", num)
	} else if num%10 >= 2 && num%10 <= 4 && (num%100 < 10 || num%100 >= 20) {
		return fmt.Sprintf("%d сообщения", num)
	}
	return fmt.Sprintf("%d сообщений", num)
}

func getRandomPhrase(phrases []string) string {
	if len(phrases) == 0 {
		return ""
	}
	return phrases[rand.Intn(len(phrases))]
}

func sendNotification(bot *tgbotapi.BotAPI, message *tgbotapi.Message, mediaGroupCache map[string]bool, link string) {
	if message.MediaGroupID != "" {
		mediaGroupCache[message.MediaGroupID] = true
	}
	messageFromUserName := message.From.UserName
	textNotification := fmt.Sprintf("Сообщение от %s удалено: обнаружен негатив. Посмотреть его можно тут: %s", messageFromUserName, link)
	msg := tgbotapi.NewMessage(message.Chat.ID, textNotification)
	//msg.ParseMode = "Markdown"
	if _, err := bot.Send(msg); err != nil {
		log.Errorf("Failed to send notification: %v", err)
	}
}
