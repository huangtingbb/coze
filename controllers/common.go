package controllers

import (
	"coze-agent-platform/utils"
	"coze-agent-platform/utils/coze"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	// 从表单获取上传的文件
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "获取文件失败: "+err.Error())
		return
	}

	// 打开上传的文件
	file, err := fileHeader.Open()
	if err != nil {
		utils.BadRequest(c, "打开文件失败: "+err.Error())
		return
	}
	defer file.Close()

	// 创建 Coze 客户端
	cozeClient, err := coze.New()
	if err != nil {
		utils.BadRequest(c, "创建Coze客户端失败: "+err.Error())
		return
	}

	// 上传文件到 Coze
	fileID, err := cozeClient.Upload(file)
	if err != nil {
		utils.BadRequest(c, "上传文件到Coze失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"file_id":  fileID,
		"filename": fileHeader.Filename,
		"size":     fileHeader.Size,
	})
}
