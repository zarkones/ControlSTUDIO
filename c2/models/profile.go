package models

import "time"

type Profile struct {
	ID                string `gorm:"primaryKey"`
	Name              string
	Description       string
	SerializedProfile string
	CreatedAt         time.Time
}
