package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Database struct {
	DB *sql.DB
}

// NewDatabase инициализирует соединение с базой данных
func NewDatabase(dsn string) (*Database, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)

	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to the database")
	return &Database{DB: db}, nil
}

// SaveMessage сохраняет сообщение в базу данных
func (d *Database) SaveMessage(userID int64, username, contentType, content string, chatID int64) error {
	query := `
		INSERT INTO messages (user_id, username, content_type, content, chat_id, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`
	_, err := d.DB.Exec(query, userID, username, contentType, content, chatID)
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}
	return nil
}

// Close закрывает соединение с базой данных
func (d *Database) Close() error {
	return d.DB.Close()
}
