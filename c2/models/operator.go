package models

import (
	"time"

	"gorm.io/gorm"
)

type Operator struct {
	Username     string `gorm:"primaryKey"`
	PublicKeyHex string
	CreatedAt    time.Time
}

func (o *Operator) BeforeCreate(tx *gorm.DB) (err error) {
	o.CreatedAt = time.Now()
	return nil
}
