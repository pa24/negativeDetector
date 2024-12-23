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
	MostActiveUser         string
	MostActiveUserMessages int
	TopVoiceUser           string
	TopVoiceMessages       int
	TopVideoUser           string
	TopVideoMessages       int
}

// GetDailyStats получает ежедневную статистику
func GetDailyStats(db *Database, chatID int64) (*DailyStats, error) {
	// Определяем текущую дату
	yesterday := time.Now().UTC().AddDate(0, 0, -1).Truncate(24 * time.Hour)

	// Определяем начало и конец вчерашнего дня
	startOfDayUTC := yesterday
	endOfDayUTC := yesterday.Add(24 * time.Hour)

	// Преобразуем время из UTC в локальное время (UTC+3)
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, fmt.Errorf("can't load location %w", err)
	}

	startOfDayLocal := startOfDayUTC.In(loc)
	endOfDayLocal := endOfDayUTC.In(loc)

	var stats DailyStats

	// 1. Общее количество сообщений
	err = db.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM messages 
		WHERE chat_id = $1 AND created_at >= $2 AND created_at < $3`,
		chatID, startOfDayLocal, endOfDayLocal,
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
		chatID, startOfDayLocal, endOfDayLocal,
	).Scan(&stats.TotalWords)
	if err != nil {
		return nil, fmt.Errorf("failed to get total words: %w", err)
	}

	// 3. Самый активный пользователь по количеству сообщений
	err = db.DB.QueryRow(`
		SELECT username, COUNT(*) 
		FROM messages
		WHERE chat_id = $1 AND created_at >= $2 AND created_at < $3
		GROUP BY username
		ORDER BY COUNT(*) DESC
		LIMIT 1`,
		chatID, startOfDayLocal, endOfDayLocal,
	).Scan(&stats.MostActiveUser, &stats.MostActiveUserMessages)
	if err != nil {
		log.Printf("failed to get most active user: %v", err)
	}

	// 4. Пользователь, отправивший больше всего голосовых сообщений
	err = db.DB.QueryRow(`
		SELECT username, COUNT(*)
		FROM messages
		WHERE chat_id = $1 AND content_type = 'voice' AND created_at >= $2 AND created_at < $3
		GROUP BY username
		ORDER BY COUNT(*) DESC
		LIMIT 1`,
		chatID, startOfDayLocal, endOfDayLocal,
	).Scan(&stats.TopVoiceUser, &stats.TopVoiceMessages)

	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			stats.TopVoiceUser = ""
			stats.TopVoiceMessages = 0
		} else {
			log.Printf("failed to get top voice user: %v", err)
		}
	}

	// 5. Пользователь, отправивший больше всего видео
	err = db.DB.QueryRow(`
		SELECT username, COUNT(*)
		FROM messages
		WHERE chat_id = $1 AND content_type = 'video_note' AND created_at >= $2 AND created_at < $3
		GROUP BY username
		ORDER BY COUNT(*) DESC
		LIMIT 1`,
		chatID, startOfDayLocal, endOfDayLocal,
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

	return &stats, nil
}
