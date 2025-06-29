package controllers

import (
	"coze-agent-platform/utils"
	"coze-agent-platform/utils/coze"

	"github.com/gin-gonic/gin"
)

// GetCozeToken 获取Coze访问令牌
// @Summary 获取Coze访问令牌
// @Description 获取Coze API的访问令牌
// @Tags Coze
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/coze/token [get]
func GetCozeToken(c *gin.Context) {
	token, err := coze.GetToken()
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, token)
}
