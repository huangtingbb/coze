package services

import (
	"coze-agent-platform/models"
	"errors"

	"gorm.io/gorm"
)

type messageService struct{}

func NewMessageService() models.MessageService {
	return &messageService{}
}

func (s *messageService) CreateMessage(message *models.Message) error {
	return models.DB.Create(message).Error
}

func (s *messageService) GetMessageById(id uint) (*models.Message, error) {
	var message models.Message
	err := models.DB.First(&message, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("消息不存在")
		}
		return nil, err
	}
	return &message, nil
}

func (s *messageService) GetMessagesByConversationId(conversationId uint, limit int) ([]*models.Message, error) {
	var messages []*models.Message
	query := models.DB.Where("conversation_id = ?", conversationId).Order("created_at ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&messages).Error
	return messages, err
}

func (s *messageService) GetRecentMessages(conversationId uint, limit int) ([]*models.Message, error) {
	var messages []*models.Message

	// 获取最近的消息，按创建时间倒序，然后取指定数量
	err := models.DB.Where("conversation_id = ?", conversationId).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	// 将消息按时间正序排列（最旧的在前面）
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (s *messageService) UpdateMessage(message *models.Message) error {
	return models.DB.Save(message).Error
}

func (s *messageService) DeleteMessage(id uint) error {
	return models.DB.Delete(&models.Message{}, id).Error
}

func (s *messageService) ListMessages(page, pageSize int) ([]*models.Message, int64, error) {
	var messages []*models.Message
	var total int64

	// 计算总数
	if err := models.DB.Model(&models.Message{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := models.DB.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&messages).Error
	return messages, total, err
}
