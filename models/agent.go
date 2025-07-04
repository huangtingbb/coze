package models

import (
	"time"

	"gorm.io/gorm"
)

type Agent struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
	Prompt      string `gorm:"type:text" json:"prompt"`
	Config      string `gorm:"type:json" json:"config"`
	Status      int    `gorm:"default:1" json:"status"` // 1:启用 0:禁用
	UserID      uint   `gorm:"not null" json:"user_id"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Agent) TableName() string {
	return "agents"
}

type AgentService interface {
	CreateAgent(agent *Agent) error
	GetAgentByID(id uint) (*Agent, error)
	GetAgentsByUserID(userId uint) ([]*Agent, error)
	UpdateAgent(agent *Agent) error
	DeleteAgent(id uint) error
	ListAgents(page, pageSize int) ([]*Agent, int64, error)
}
