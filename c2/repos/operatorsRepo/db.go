package operatorsRepo

import (
	"c2/db"
	"c2/models"
)

func Get(userID string) (operator models.Operator, err error) {
	return operator, db.ORM.Where("user_id = ?", userID).First(&operator).Error
}

func GetMultiple() (operators []models.Operator, err error) {
	return operators, db.ORM.Find(&operators).Error
}

func Insert(operator *models.Operator) (err error) {
	return db.ORM.Create(&operator).Error
}
