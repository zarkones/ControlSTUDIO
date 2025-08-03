package models

import "time"

type Agent struct {
	ID        string `gorm:"primaryKey"`
	IP        string
	OS        string
	Arch      string
	Hostname  string
	CreatedAt time.Time
}
