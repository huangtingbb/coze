package models

import (
	"time"

	"gorm.io/gorm"
)

type Conversation struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Title   string `json:"title"`
	Status  int    `gorm:"default:1" json:"status"` // 1:进行中 0:已结束
	UserID  uint   `gorm:"not null" json:"user_id"`
	AgentID uint   `gorm:"not null" json:"agent_id"`

	// 关联关系
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Agent    Agent     `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	Messages []Message `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
}

func (Conversation) TableName() string {
	return "conversations"
}

type ConversationService interface {
	CreateConversation(conversation *Conversation) error
	GetConversationByID(id uint) (*Conversation, error)
	GetConversationsByUserID(userID uint) ([]*Conversation, error)
	UpdateConversation(conversation *Conversation) error
	DeleteConversation(id uint) error
	ListConversations(page, pageSize int) ([]*Conversation, int64, error)
}
