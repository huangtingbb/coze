package controllers

import (
	"coze-agent-platform/services"
	"coze-agent-platform/utils"

	"github.com/gin-gonic/gin"
)

// GetUserProfile 获取用户资料
// @Summary 获取用户资料
// @Description 获取当前登录用户的资料信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /users/profile [get]
func GetUserProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	userService := services.NewUserService()
	user, err := userService.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, user)
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// UpdateUserProfile 更新用户资料
// @Summary 更新用户资料
// @Description 更新当前登录用户的资料信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body UpdateProfileRequest true "用户资料"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /users/profile [put]
func UpdateUserProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误："+err.Error())
		return
	}

	userService := services.NewUserService()
	user, err := userService.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	// 更新用户信息
	user.Nickname = req.Nickname
	user.Avatar = req.Avatar

	if err := userService.UpdateUser(user); err != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", user)
}
