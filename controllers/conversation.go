package controllers

import (
	"coze-agent-platform/utils"
	"coze-agent-platform/utils/coze"

	"github.com/gin-gonic/gin"
)

type CreateConversationRequest struct {
	AgentID uint   `json:"agent_id" binding:"required"`
	Title   string `json:"title"`
}

type SendMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

// ListConversations 获取对话列表
// @Summary 获取对话列表
// @Description 获取当前用户的对话列表
// @Tags 对话
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} utils.PageResponse
// @Failure 401 {object} utils.Response
// @Router /api/conversations [get]
func ListConversations(c *gin.Context) {
	utils.Success(c, []interface{}{})
}

// CreateConversation 创建对话
// @Summary 创建对话
// @Description 创建新的对话会话
// @Tags 对话
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreateConversationRequest true "对话信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/conversations [post]
func CreateConversation(c *gin.Context) {
	cozeConv, err := coze.NewConversation()
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	id, err := cozeConv.CreateConversation()
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, id)
}

// GetConversation 获取对话详情
// @Summary 获取对话详情
// @Description 根据ID获取对话详细信息
// @Tags 对话
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "对话ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/conversations/{id} [get]
func GetConversation(c *gin.Context) {
	utils.SuccessWithMessage(c, "功能开发中", nil)
}

// DeleteConversation 删除对话
// @Summary 删除对话
// @Description 删除指定的对话
// @Tags 对话
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "对话ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/conversations/{id} [delete]
func DeleteConversation(c *gin.Context) {
	utils.SuccessWithMessage(c, "功能开发中", nil)
}

// GetMessages 获取消息列表
// @Summary 获取消息列表
// @Description 获取指定对话的消息列表
// @Tags 消息
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "对话ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} utils.PageResponse
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/conversations/{id}/messages [get]
func GetMessages(c *gin.Context) {
	utils.Success(c, []interface{}{})
}

// SendMessage 发送消息
// @Summary 发送消息
// @Description 向指定对话发送消息
// @Tags 消息
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "对话ID"
// @Param request body SendMessageRequest true "消息内容"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/conversations/{id}/messages [post]
func SendMessage(c *gin.Context) {
	conversationID := c.Param("id")
	content := c.PostForm("content")
	cozeConv, err := coze.NewConversation()
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	msgs,err := cozeConv.SendMessage(conversationID, content)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c,msgs)
}
