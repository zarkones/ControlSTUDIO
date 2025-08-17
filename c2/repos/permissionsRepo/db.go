package permissionsRepo

import (
	"c2/db"

	access "github.com/zarkones/ControlACCESS"
)

func Get(userID string, permissionKey access.PermissionKey) (permission access.Permission, err error) {
	return permission, db.ORM.
		Where("user_id = ?", userID).
		Where("permission_key = ?", permissionKey).
		First(&permission).Error
}

func GetMultipleByUserID(userID string) (permissions []access.Permission, err error) {
	return permissions, db.ORM.Where("user_id = ?", userID).Find(&permissions).Error
}

func GetMultiple() (permissions []access.Permission, err error) {
	return permissions, db.ORM.Find(&permissions).Error
}

func Insert(permission *access.Permission) (err error) {
	return db.ORM.Create(&permission).Error
}

func Delete(userID string, permissionKey access.PermissionKey) (err error) {
	return db.ORM.
		Where("user_id = ?", userID).
		Where("permission_key = ?", permissionKey).
		Delete(access.Permission{}).Error
}
