package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID        string `gorm:"primaryKey"`
	AgentID   string
	Request   string
	Response  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	m.CreatedAt = time.Now()
	return nil
}
