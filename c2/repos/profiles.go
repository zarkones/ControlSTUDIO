package repos

import (
	"c2/db"
	"c2/models"
)

func GetProfile(profileID string) (profile models.MetaProfile, err error) {
	return profile, db.ORM.Where("profile_id = ?", profileID).First(&profile).Error
}

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
