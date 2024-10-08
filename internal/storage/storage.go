package storage

import (
	"schedule-bot/pkg/models"
)

type Storage interface {
	Add(message *models.Message) error
	Fetch(chatID int64) (*models.Message, error)
	Update(msg *models.Message) error
}
