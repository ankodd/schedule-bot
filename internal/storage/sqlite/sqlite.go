package sqlite

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"schedule-bot/pkg/models"
)

const (
	MessageNotFound = "message not found"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) (*Storage, error) {
	q := `CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY,
		chatID INTEGER NOT NULL,
		messageID INTEGER NOT NULL,
		filename VARCHAR NOT NULL
	)`

	_, err := db.Exec(q)
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Add(message *models.Message) error {
	q := `INSERT INTO messages (chatID, messageID, filename) VALUES (?, ?, ?)`
	_, err := s.db.Exec(q, message.ChatID, message.MessageID, message.Filename)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Fetch(chatID int64) (*models.Message, error) {
	q := `SELECT * FROM messages WHERE chatID = ?`
	row := s.db.QueryRow(q, chatID)

	msg := models.Message{}

	err := row.Scan(&msg.ID, &msg.ChatID, &msg.MessageID, &msg.Filename)
	if err != nil && err.Error() == sql.ErrNoRows.Error() {
		return nil, errors.New(MessageNotFound)
	} else if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (s *Storage) Update(msg *models.Message) error {
	q := `UPDATE messages SET messageID = ? WHERE chatID = ?`
	_, err := s.db.Exec(q, msg.MessageID, msg.ChatID)
	if err != nil {
		return err
	}

	return nil
}
