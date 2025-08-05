package repos

import (
	"c2/db"
	"c2/models"
	"errors"
	"time"
)

var ErrMsgRespPopulated = errors.New("message's response is already populated")

func GetMessages(agentID string, offset, limit int) (messages []models.Message, err error) {
	return messages, db.ORM.
		Where("agent_id = ?", agentID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&messages).Error
}

func GetMessagesBefore(agentID string, before int64, limit int) (messages []models.Message, err error) {
	return messages, db.ORM.
		Where("agent_id = ?", agentID).
		Where("created_at < ?", before).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
}

func GetMessagesAfter(agentID string, after int64, limit int) (messages []models.Message, err error) {
	return messages, db.ORM.
		Where("agent_id = ?", agentID).
		Where("created_at > ?", after).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error
}

// func Get(messageID string) (messages models.Message, err error) {
// 	return messages, db.ORM.Where("id = ?", messageID).First(&messages).Error
// }

// func GetMultipleByIDs(messageIDs []string) (messages []models.Message, err error) {
// 	return messages, db.ORM.
// 		Find(&messages, messageIDs).Error
// }

func GetOldestMessageForAgent(agentID string) (message models.Message, err error) {
	return message, db.ORM.
		Where("agent_id = ?", agentID).
		Where("response = ?", "").
		Order("created_at ASC").
		First(&message).Error
}

// func GetMultipleForAgent(agentID string) (messages []models.Message, err error) {
// 	return messages, db.ORM.Where("agent_id = ?", agentID).Where("response = ?", "").Find(&messages).Error
// }

func InsertMessage(message *models.Message) (err error) {
	return db.ORM.Create(&message).Error
}

func UpdateOldestMessageResponse(agentID, response string) (err error) {
	message, err := GetOldestMessageForAgent(agentID)
	if err != nil {
		return err
	}
	if message.Response != "" {
		return ErrMsgRespPopulated
	}
	message.Response = response
	message.UpdatedAt = time.Now().UnixNano()
	return db.ORM.Save(message).Error
}
