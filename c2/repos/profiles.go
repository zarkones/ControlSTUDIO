package repos

import (
	"c2/db"
	"c2/models"
)

func GetProfiles() (profiles []models.MetaProfile, err error) {
	return profiles, db.ORM.Find(&profiles).Error
}

func InsertProfile(profile *models.MetaProfile) (err error) {
	return db.ORM.Create(profile).Error
}

func UpsertProfile(profile *models.MetaProfile) (err error) {
	return db.ORM.Save(profile).Error
}

func DeleteProfile(profileID string) (err error) {
	return db.ORM.Delete(&models.MetaProfile{ID: profileID}).Error
}
