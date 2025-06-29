package controllers

import (
	"coze-agent-platform/models"
	"coze-agent-platform/services"
	"coze-agent-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateAgentRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
	Prompt      string `json:"prompt"`
	Config      string `json:"config"`
}

// CreateAgent 创建Agent
// @Summary 创建Agent
// @Description 创建新的AI Agent
// @Tags Agent
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreateAgentRequest true "Agent信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/agents [post]
func CreateAgent(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误："+err.Error())
		return
	}

	agent := &models.Agent{
		Name:        req.Name,
		Description: req.Description,
		Avatar:      req.Avatar,
		Prompt:      req.Prompt,
		Config:      req.Config,
		Status:      1,
		UserID:      userID,
	}

	agentService := services.NewAgentService()
	if err := agentService.CreateAgent(agent); err != nil {
		utils.InternalServerError(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", agent)
}

// GetAgent 获取Agent详情
// @Summary 获取Agent详情
// @Description 根据ID获取Agent详细信息
// @Tags Agent
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Agent ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/agents/{id} [get]
func GetAgent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID格式错误")
		return
	}

	agentService := services.NewAgentService()
	agent, err := agentService.GetAgentByID(uint(id))
	if err != nil {
		utils.NotFound(c, "Agent不存在")
		return
	}

	utils.Success(c, agent)
}

// ListAgents 获取Agent列表
// @Summary 获取Agent列表
// @Description 获取当前用户的Agent列表
// @Tags Agent
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} utils.PageResponse
// @Failure 401 {object} utils.Response
// @Router /api/agents [get]
func ListAgents(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	agentService := services.NewAgentService()
	agents, total, err := agentService.ListAgents(page, size)
	if err != nil {
		utils.InternalServerError(c, "查询失败")
		return
	}

	// 过滤当前用户的Agent
	var userAgents []*models.Agent
	for _, agent := range agents {
		if agent.UserID == userID {
			userAgents = append(userAgents, agent)
		}
	}

	utils.PageSuccess(c, userAgents, total, page, size)
}

// UpdateAgent 更新Agent
// @Summary 更新Agent
// @Description 更新Agent信息
// @Tags Agent
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Agent ID"
// @Param request body CreateAgentRequest true "Agent信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/agents/{id} [put]
func UpdateAgent(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID格式错误")
		return
	}

	var req CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误："+err.Error())
		return
	}

	agentService := services.NewAgentService()
	agent, err := agentService.GetAgentByID(uint(id))
	if err != nil {
		utils.NotFound(c, "Agent不存在")
		return
	}

	// 检查权限
	if agent.UserID != userID {
		utils.Unauthorized(c, "无权限操作")
		return
	}

	// 更新Agent信息
	agent.Name = req.Name
	agent.Description = req.Description
	agent.Avatar = req.Avatar
	agent.Prompt = req.Prompt
	agent.Config = req.Config

	if err := agentService.UpdateAgent(agent); err != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", agent)
}

// DeleteAgent 删除Agent
// @Summary 删除Agent
// @Description 删除指定的Agent
// @Tags Agent
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Agent ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/agents/{id} [delete]
func DeleteAgent(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID格式错误")
		return
	}

	agentService := services.NewAgentService()
	agent, err := agentService.GetAgentByID(uint(id))
	if err != nil {
		utils.NotFound(c, "Agent不存在")
		return
	}

	// 检查权限
	if agent.UserID != userID {
		utils.Unauthorized(c, "无权限操作")
		return
	}

	if err := agentService.DeleteAgent(uint(id)); err != nil {
		utils.InternalServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
