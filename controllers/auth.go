package controllers

import (
	"coze-agent-platform/config"
	"coze-agent-platform/models"
	"coze-agent-platform/services"
	"coze-agent-platform/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取JWT token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} utils.Response{data=AuthResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误："+err.Error())
		return
	}

	userService := services.NewUserService()
	user, err := userService.GetUserByUsername(req.Username)
	if err != nil {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, config.Cfg.JWT.Secret, config.Cfg.JWT.Expire)
	if err != nil {
		utils.InternalServerError(c, "Token生成失败")
		return
	}

	utils.Success(c, AuthResponse{
		Token: token,
		User:  *user,
	})
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册新账户
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 200 {object} utils.Response{data=AuthResponse}
// @Failure 400 {object} utils.Response
// @Router /api/auth/register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误："+err.Error())
		return
	}

	userService := services.NewUserService()

	// 检查用户名是否已存在
	if _, err := userService.GetUserByUsername(req.Username); err == nil {
		utils.BadRequest(c, "用户名已存在")
		return
	}

	// 检查邮箱是否已存在
	if _, err := userService.GetUserByEmail(req.Email); err == nil {
		utils.BadRequest(c, "邮箱已存在")
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	// 创建用户
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Status:   1,
		Role:     1,
	}

	if err := userService.CreateUser(user); err != nil {
		utils.InternalServerError(c, "用户创建失败")
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, config.Cfg.JWT.Secret, config.Cfg.JWT.Expire)
	if err != nil {
		utils.InternalServerError(c, "Token生成失败")
		return
	}

	utils.Success(c, AuthResponse{
		Token: token,
		User:  *user,
	})
}
