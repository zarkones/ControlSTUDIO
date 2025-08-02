package models

import "time"

type Agent struct {
	ID        string `gorm:"primaryKey"`
	IP        string
	CreatedAt time.Time
}
