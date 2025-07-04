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

	CozeConversationID string `gorm:"column:coze_conversation_id;size:100;not null" json:"coze_conversation_id"`
	UserId             uint   `gorm:"column:user_id;not null" json:"user_id"`
	Title              string `gorm:"column:title;size:100" json:"title"`

	// 关联关系
	User     User      `gorm:"foreignKey:UserId" json:"user,omitempty"`
	Messages []Message `gorm:"foreignKey:ConversationId" json:"messages,omitempty"`
}

func (Conversation) TableName() string {
	return "conversation"
}

type ConversationService interface {
	CreateConversation(conversation *Conversation) error
	GetConversationById(id uint) (*Conversation, error)
	GetConversationByCozeId(cozeConversationId string) (*Conversation, error)
	GetConversationsByUserId(userId uint) ([]*Conversation, error)
	UpdateConversation(conversation *Conversation) error
	DeleteConversation(id uint) error
	ListConversations(page, pageSize int) ([]*Conversation, int64, error)
}
