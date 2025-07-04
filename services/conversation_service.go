package services

import (
	"coze-agent-platform/models"
	"coze-agent-platform/utils/coze"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type conversationService struct{}

func NewConversationService() models.ConversationService {
	return &conversationService{}
}

func (s *conversationService) CreateConversation(conversation *models.Conversation) error {
	// 创建Coze对话
	cozeConv, err := coze.New()
	if err != nil {
		return fmt.Errorf("初始化Coze对话失败: %v", err.Error())
	}

	cozeConversationID, err := cozeConv.CreateConversation()
	if err != nil {
		return fmt.Errorf("创建Coze对话失败: %v", err.Error())
	}

	conversation.CozeConversationID = cozeConversationID
	return models.DB.Create(conversation).Error
}

func (s *conversationService) GetConversationById(id uint) (*models.Conversation, error) {
	var conversation models.Conversation
	err := models.DB.First(&conversation, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("对话不存在")
		}
		return nil, err
	}
	return &conversation, nil
}

func (s *conversationService) GetConversationByCozeId(cozeConversationId string) (*models.Conversation, error) {
	var conversation models.Conversation
	err := models.DB.Where("coze_conversation_id = ?", cozeConversationId).First(&conversation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("对话不存在")
		}
		return nil, err
	}
	return &conversation, nil
}

func (s *conversationService) GetConversationsByUserId(userId uint) ([]*models.Conversation, error) {
	var conversations []*models.Conversation
	err := models.DB.Where("user_id = ?", userId).Order("created_at DESC").Find(&conversations).Error
	return conversations, err
}

func (s *conversationService) UpdateConversation(conversation *models.Conversation) error {
	return models.DB.Save(conversation).Error
}

func (s *conversationService) DeleteConversation(id uint) error {
	return models.DB.Delete(&models.Conversation{}, id).Error
}

func (s *conversationService) ListConversations(page, pageSize int) ([]*models.Conversation, int64, error) {
	var conversations []*models.Conversation
	var total int64

	// 计算总数
	if err := models.DB.Model(&models.Conversation{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := models.DB.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&conversations).Error
	return conversations, total, err
}
