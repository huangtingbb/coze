package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	CozeMessageId  string `gorm:"column:coze_message_id;size:100;not null" json:"coze_message_id"`
	ConversationId uint   `gorm:"column:conversation_id;not null" json:"conversation_id"`
	ModelId        uint   `gorm:"column:model_id;not null" json:"model_id"`
	Metadata       string `gorm:"column:metadata;size:255" json:"metadata"`
	Role           string `gorm:"column:role;size:10;not null" json:"role"` // user、assistant、system
	Content        string `gorm:"column:content;type:text;not null" json:"content"`
	Tokens         int    `gorm:"column:tokens;default:0" json:"tokens"`

	// 关联关系
	Conversation Conversation `gorm:"foreignKey:ConversationId" json:"conversation,omitempty"`
}

func (Message) TableName() string {
	return "message"
}

type MessageService interface {
	CreateMessage(message *Message) error
	GetMessageById(id uint) (*Message, error)
	GetMessagesByConversationId(conversationId uint, limit int) ([]*Message, error)
	GetRecentMessages(conversationId uint, limit int) ([]*Message, error)
	UpdateMessage(message *Message) error
	DeleteMessage(id uint) error
	ListMessages(page, pageSize int) ([]*Message, int64, error)
}
