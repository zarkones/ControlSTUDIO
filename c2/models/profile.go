package models

import (
	"common/profiles"
	"encoding/json"
	"time"
)

type MetaProfile struct {
	ID                string `gorm:"primaryKey"`
	Name              string
	Description       string
	SerializedProfile string
	CreatedAt         time.Time
}

func (p *MetaProfile) GetProfile() (profile profiles.Profile, err error) {
	return profile, json.Unmarshal([]byte(p.SerializedProfile), &profile)
}
