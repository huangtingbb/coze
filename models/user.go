package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Status   int    `gorm:"default:1" json:"status"` // 1:正常 0:禁用
	Role     int    `gorm:"default:1" json:"role"`   // 1:普通用户 2:管理员
}

func (User) TableName() string {
	return "users"
}

// UserService 用户服务接口
type UserService interface {
	CreateUser(user *User) error
	GetUserByID(id uint) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uint) error
	ListUsers(page, pageSize int) ([]*User, int64, error)
}
