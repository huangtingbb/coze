package services

import (
	"coze-agent-platform/models"
	"errors"

	"gorm.io/gorm"
)

type userService struct{}

func NewUserService() models.UserService {
	return &userService{}
}

func (s *userService) CreateUser(user *models.User) error {
	return models.DB.Create(user).Error
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := models.DB.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := models.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := models.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

func (s *userService) UpdateUser(user *models.User) error {
	return models.DB.Save(user).Error
}

func (s *userService) DeleteUser(id uint) error {
	return models.DB.Delete(&models.User{}, id).Error
}

func (s *userService) ListUsers(page, pageSize int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	// 计算总数
	if err := models.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := models.DB.Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}
