package coze

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/coze-dev/coze-go"
)

// FileUploadResponse 文件上传响应结构

// Upload 上传文件到 Coze
func (client *Client) Upload(file io.Reader) (string, error) {
	// 创建临时文件
	tempFile, err := os.CreateTemp("", "coze_upload_*")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tempFile.Name()) // 清理临时文件

	// 将文件内容复制到临时文件
	_, err = io.Copy(tempFile, file)
	if err != nil {
		tempFile.Close()
		return "", fmt.Errorf("复制文件内容失败: %v", err)
	}

	// 重置文件指针到开始位置
	tempFile.Seek(0, io.SeekStart)

	uploadReq := &coze.UploadFilesReq{
		File: tempFile,
	}

	uploadResp, err := client.Api.Files.Upload(context.Background(), uploadReq)
	tempFile.Close() // 关闭临时文件
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %v", err)
	}

	return uploadResp.FileInfo.ID, nil
}
