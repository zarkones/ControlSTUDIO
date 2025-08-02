package repos

import (
	"c2/db"
	"c2/models"
)

func GetAgents() (agents []models.Agent, err error) {
	return agents, db.ORM.Find(&agents).Error
}

func InsertAgent(agent *models.Agent) (err error) {
	return db.ORM.Create(agent).Error
}
