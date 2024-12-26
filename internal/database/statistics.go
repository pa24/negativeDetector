package database

import (
	"database/sql"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

type DailyStats struct {
	TotalMessages          int
	TotalWords             int
	MostActiveUserID       string
	MostActiveUserMessages int
	TopVoiceUser           string
	TopVoiceMessages       int
	TopVideoUser           string
	TopVideoMessages       int
}

// GetDailyStats получает ежедневную статистику
func GetDailyStats(db *Database, chatID int64) (*DailyStats, error) {
	// Преобразуем время из UTC в локальное время (UTC+3)
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {

	}
	now := time.Now().In(loc)

	// Определяем вчерашний день
	yesterdayLocal := now.AddDate(0, 0, -1)
	startOfDayLocal := time.Date(yesterdayLocal.Year(), yesterdayLocal.Month(), yesterdayLocal.Day(), 0, 0, 0, 0, loc)
	endOfDayLocal := startOfDayLocal.Add(24 * time.Hour)

	startOfDayUTC := startOfDayLocal.UTC()
	endOfDayUTC := endOfDayLocal.UTC()

	var stats DailyStats

	userIDToName := map[string]string{
		"1633921608": "Вера",
		"1302900055": "Аня",
		"927178234":  "Андрей",
		"1701035449": "Славик",
		"1705013271": "Гоша",
		"919183144":  "Коля",
		"5136453150": "Яна",
	}

	// 1. Общее количество сообщений
	err = db.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM messages 
		WHERE chat_id = $1 AND created_at >= $2 AND created_at < $3`,
		chatID, startOfDayUTC, endOfDayUTC,
	).Scan(&stats.TotalMessages)
	if err != nil {
		return nil, fmt.Errorf("failed to get total messages: %w", err)
	}

	// 2. Общее количество слов
	err = db.DB.QueryRow(`
		SELECT COALESCE(SUM(LENGTH(content) - LENGTH(REPLACE(content, ' ', '')) + 1), 0)
		FROM messages
		WHERE chat_id = $1 AND content_type = 'text' AND created_at >= $2 AND created_at < $3
		AND content <> '<empty>'`,
		chatID, startOfDayUTC, endOfDayUTC,
	).Scan(&stats.TotalWords)
	if err != nil {
		return nil, fmt.Errorf("failed to get total words: %w", err)
	}

	// 3. Самый активный пользователь по количеству сообщений
	err = db.DB.QueryRow(`
		SELECT user_id, COUNT(*) 
		FROM messages
		WHERE chat_id = $1 AND created_at >= $2 AND created_at < $3
		GROUP BY user_id
		ORDER BY COUNT(*) DESC
		LIMIT 1`,
		chatID, startOfDayUTC, endOfDayUTC,
	).Scan(&stats.MostActiveUserID, &stats.MostActiveUserMessages)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			stats.MostActiveUserID = ""
			stats.MostActiveUserMessages = 0
		} else {
			log.Printf("failed to get top voice user: %v", err)
		}
	}

	mostActiveUserName, ok := userIDToName[stats.MostActiveUserID]
	if ok {
		stats.MostActiveUserID = mostActiveUserName
	} else {
		stats.MostActiveUserID = "никто" // Если имя пользователя не найдено
	}

	// 4. Пользователь, отправивший больше всего голосовых сообщений
	err = db.DB.QueryRow(`
		SELECT user_id, COUNT(*)
		FROM messages
		WHERE chat_id = $1 AND content_type = 'voice' AND created_at >= $2 AND created_at < $3
		GROUP BY user_id
		ORDER BY COUNT(*) DESC
		LIMIT 1`,
		chatID, startOfDayUTC, endOfDayUTC,
	).Scan(&stats.TopVoiceUser, &stats.TopVoiceMessages)

	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			stats.TopVoiceUser = ""
			stats.TopVoiceMessages = 0
		} else {
			log.Printf("failed to get top voice user: %v", err)
		}
	}

	topVoiceUserName, ok := userIDToName[stats.TopVoiceUser]
	if ok {
		stats.TopVoiceUser = topVoiceUserName
	} else {
		stats.TopVoiceUser = "никто" // Если имя пользователя не найдено
	}

	// 5. Пользователь, отправивший больше всего видео
	err = db.DB.QueryRow(`
		SELECT user_id, COUNT(*)
		FROM messages
		WHERE chat_id = $1 AND content_type = 'video_note' AND created_at >= $2 AND created_at < $3
		GROUP BY user_id
		ORDER BY COUNT(*) DESC
		LIMIT 1`,
		chatID, startOfDayUTC, endOfDayUTC,
	).Scan(&stats.TopVideoUser, &stats.TopVideoMessages)

	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			// Handle case where no rows are returned
			stats.TopVideoUser = ""
			stats.TopVideoMessages = 0
		} else {
			log.Printf("failed to get top video user: %v", err)
		}
	}
	topVideoUserName, ok := userIDToName[stats.TopVideoUser]
	if ok {
		stats.TopVideoUser = topVideoUserName
	} else {
		stats.TopVideoUser = "никто" // Если имя пользователя не найдено
	}

	return &stats, nil
}
	