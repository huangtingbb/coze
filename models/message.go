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

	Content        string `gorm:"type:text;not null" json:"content"`
	MessageType    int    `gorm:"not null" json:"message_type"` // 1:用户消息 2:AI回复
	ConversationID uint   `gorm:"not null" json:"conversation_id"`
	UserID         uint   `gorm:"not null" json:"user_id"`

	// 关联关系
	Conversation Conversation `gorm:"foreignKey:ConversationID" json:"conversation,omitempty"`
	User         User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Message) TableName() string {
	return "messages"
}

type MessageService interface {
	CreateMessage(message *Message) error
	GetMessageByID(id uint) (*Message, error)
	GetMessagesByConversationID(conversationID uint) ([]*Message, error)
	UpdateMessage(message *Message) error
	DeleteMessage(id uint) error
	ListMessages(page, pageSize int) ([]*Message, int64, error)
}
