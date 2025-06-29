package services

import (
	"coze-agent-platform/models"
	"errors"

	"gorm.io/gorm"
)

type agentService struct{}

func NewAgentService() models.AgentService {
	return &agentService{}
}

func (s *agentService) CreateAgent(agent *models.Agent) error {
	return models.DB.Create(agent).Error
}

func (s *agentService) GetAgentByID(id uint) (*models.Agent, error) {
	var agent models.Agent
	err := models.DB.Preload("User").First(&agent, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Agent不存在")
		}
		return nil, err
	}
	return &agent, nil
}

func (s *agentService) GetAgentsByUserID(userID uint) ([]*models.Agent, error) {
	var agents []*models.Agent
	err := models.DB.Where("user_id = ?", userID).Find(&agents).Error
	return agents, err
}

func (s *agentService) UpdateAgent(agent *models.Agent) error {
	return models.DB.Save(agent).Error
}

func (s *agentService) DeleteAgent(id uint) error {
	return models.DB.Delete(&models.Agent{}, id).Error
}

func (s *agentService) ListAgents(page, pageSize int) ([]*models.Agent, int64, error) {
	var agents []*models.Agent
	var total int64

	// 计算总数
	if err := models.DB.Model(&models.Agent{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := models.DB.Preload("User").Offset(offset).Limit(pageSize).Find(&agents).Error
	return agents, total, err
}
