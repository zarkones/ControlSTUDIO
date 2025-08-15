package models

import (
	"encoding/json"
	"time"

	profiles "github.com/zarkones/ControlPROFILE"
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
