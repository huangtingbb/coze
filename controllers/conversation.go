package controllers

import (
	"coze-agent-platform/models"
	"coze-agent-platform/services"
	"coze-agent-platform/utils"
	"coze-agent-platform/utils/coze"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type CreateConversationRequest struct {
	Title string `json:"title"`
}

type SendMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

var (
	conversationService = services.NewConversationService()
	messageService      = services.NewMessageService()
)

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
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "用户未登录")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	// 获取用户的对话列表
	conversations, err := conversationService.GetConversationsByUserId(userID.(uint))
	if err != nil {
		utils.InternalServerError(c, "获取对话列表失败: "+err.Error())
		return
	}

	// 简单分页
	total := len(conversations)
	start := (page - 1) * size
	end := start + size

	if start >= total {
		conversations = []*models.Conversation{}
	} else {
		if end > total {
			end = total
		}
		conversations = conversations[start:end]
	}

	utils.PageSuccess(c, conversations, int64(total), page, size)
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
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "用户未登录")
		return
	}

	var req CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数格式错误: "+err.Error())
		return
	}

	// 创建数据库记录
	conversation := &models.Conversation{
		UserId: userID.(uint),
		Title:  req.Title,
	}

	if err := conversationService.CreateConversation(conversation); err != nil {
		utils.InternalServerError(c, "保存对话失败: "+err.Error())
		return
	}

	utils.Success(c, conversation)
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
	conversationId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "对话ID格式错误")
		return
	}

	conversation, err := conversationService.GetConversationById(uint(conversationId))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, conversation)
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
	conversationId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "对话ID格式错误")
		return
	}

	if err := conversationService.DeleteConversation(uint(conversationId)); err != nil {
		utils.InternalServerError(c, "删除对话失败: "+err.Error())
		return
	}

	utils.Success(c, nil)
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
	conversationId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "对话ID格式错误")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}

	// 获取消息列表
	messages, err := messageService.GetMessagesByConversationId(uint(conversationId), size)
	if err != nil {
		utils.InternalServerError(c, "获取消息列表失败: "+err.Error())
		return
	}

	utils.PageSuccess(c, messages, int64(len(messages)), page, size)
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
	conversationIdStr := c.DefaultQuery("conversation_id", "0")
	if conversationIdStr == "" {
		conversationIdStr = "0"
	}

	conversationId, err := strconv.ParseUint(conversationIdStr, 10, 32)
	if err != nil {
		// 如果解析失败，默认为0
		conversationId = 0
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数格式错误: "+err.Error())
		return
	}

	// 获取对话信息
	conversation, err := conversationService.GetConversationById(uint(conversationId))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	// 获取最近20条历史消息
	historyMessages, err := messageService.GetRecentMessages(uint(conversationId), 20)
	if err != nil {
		utils.InternalServerError(c, "获取历史消息失败: "+err.Error())
		return
	}

	// 发送消息到Coze
	cozeConv, err := coze.New()
	if err != nil {
		utils.BadRequest(c, "初始化Coze对话失败: "+err.Error())
		return
	}

	cozeResp, err := cozeConv.SendMessage(conversation.CozeConversationID, req.Content)
	if err != nil {
		utils.BadRequest(c, "发送消息失败: "+err.Error())
		return
	}

	// 保存用户消息到数据库
	userMessage := &models.Message{
		CozeMessageId:  generateMessageId(),
		ConversationId: uint(conversationId),
		ModelId:        1,
		Role:           "user",
		Content:        req.Content,
		Tokens:         0,
	}

	if err := messageService.CreateMessage(userMessage); err != nil {
		utils.InternalServerError(c, "保存用户消息失败: "+err.Error())
		return
	}

	// 保存AI回复到数据库
	if cozeResp.Message.ID != "" {
		aiMessage := &models.Message{
			CozeMessageId:  cozeResp.Message.ID,
			ConversationId: uint(conversationId),
			ModelId:        1,
			Role:           "assistant",
			Content:        getMessageContent(cozeResp.Message),
			Tokens:         0,
		}

		if err := messageService.CreateMessage(aiMessage); err != nil {
			utils.InternalServerError(c, "保存AI回复失败: "+err.Error())
			return
		}
	}

	// 构建返回数据
	responseData := map[string]interface{}{
		"response":      cozeResp,
		"history_count": len(historyMessages),
		"user_message":  userMessage,
	}

	utils.Success(c, responseData)
}

// SendMessageStream 发送消息(流式)
// @Summary 发送消息(流式)
// @Description 向指定对话发送消息，使用 SSE 协议返回流式响应
// @Tags 消息
// @Accept json
// @Produce text/event-stream
// @Security ApiKeyAuth
// @Param id path string true "对话ID"
// @Param user_id query string true "用户ID"
// @Param request body SendMessageRequest true "消息内容"
// @Success 200 {string} string "SSE 流式响应"
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/conversations/messages/stream [post]
func SendMessageStream(c *gin.Context) {
	conversationIdStr := c.DefaultQuery("conversation_id", "0")
	if conversationIdStr == "" {
		conversationIdStr = "0"
	}

	conversationId, err := strconv.ParseUint(conversationIdStr, 10, 32)
	if err != nil {
		// 如果解析失败，默认为0
		conversationId = 0
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数格式错误: "+err.Error())
		return
	}

	userIdStr := c.DefaultQuery("user_id", "0")
	if userIdStr == "" {
		userIdStr = "0"
	}

	userID, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		// 如果解析失败，默认为0
		userID = 0
	}

	var conversation *models.Conversation
	var historyMessageList []*models.Message

	if conversationId == 0 {
		conversation = &models.Conversation{}
		conversation.UserId = uint(userID)
		conversation.Title = strings.SplitN(req.Content, "\n", 2)[0]
		err = conversationService.CreateConversation(conversation)
		if err != nil {
			utils.InternalServerError(c, "创建对话失败: "+err.Error())
			return
		}
	} else {
		// 获取对话信息
		conversation, err = conversationService.GetConversationById(uint(conversationId))
		if err != nil {
			utils.NotFound(c, err.Error())
			return
		}

		fmt.Println(conversation.ID)
		// 获取最近20条历史消息
		historyMessageList, err = messageService.GetRecentMessages(conversation.ID, 20)
		if err != nil {
			utils.InternalServerError(c, "获取历史消息失败: "+err.Error())
			return
		}
	}

	fmt.Println(historyMessageList)

	// 设置 SSE 头部
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	// 检查客户端是否断开连接
	clientGone := c.Request.Context().Done()

	cozeConv, err := coze.New()
	if err != nil {
		// SSE 错误格式
		c.SSEvent("error", map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		})
		c.Writer.Flush()
		return
	}

	// 先保存用户消息
	userMessage := &models.Message{
		CozeMessageId:  generateMessageId(),
		ConversationId: conversation.ID,
		ModelId:        1,
		Role:           "user",
		Content:        req.Content,
		Tokens:         0,
	}

	if err := messageService.CreateMessage(userMessage); err != nil {
		c.SSEvent("error", map[string]interface{}{
			"error":   true,
			"message": "保存用户消息失败: " + err.Error(),
		})
		c.Writer.Flush()
		return
	}

	var aiMessageContent strings.Builder
	var aiMessageId string
	var aiMessageTokens int
	historyMessageList = append(historyMessageList, userMessage)
	fmt.Println(historyMessageList)

	// 定义流式回调函数
	onMessage := func(eventType string, data interface{}) {
		select {
		case <-clientGone:
			return // 客户端已断开连接
		default:
			eventData := map[string]interface{}{
				"type": eventType,
				"data": data,
			}

			// 处理消息增量更新
			if eventType == "message_delta" {
				if msgData, ok := data.(map[string]interface{}); ok {
					if content, ok := msgData["content"].(string); ok {
						aiMessageContent.WriteString(content)
					}
				}
			}

			// 处理对话完成
			if eventType == "chat_completed" {
				if msgData, ok := data.(map[string]interface{}); ok {
					if chatId, ok := msgData["chat_id"].(string); ok {
						aiMessageId = chatId
					}

					if usage, ok := msgData["usage"].(map[string]interface{}); ok {
						for _, token := range usage {
							aiMessageTokens += token.(int)
						}
					}
				}

				// 保存AI回复消息
				if aiMessageContent.Len() > 0 {
					aiMessage := &models.Message{
						CozeMessageId:  aiMessageId,
						ConversationId: conversation.ID,
						ModelId:        1,
						Role:           "assistant",
						Content:        aiMessageContent.String(),
						Tokens:         aiMessageTokens,
					}

					if err := messageService.CreateMessage(aiMessage); err != nil {
						// 错误处理，但不中断流式响应
						fmt.Printf("保存AI回复失败: %v\n", err)
					}
				}
			}

			c.SSEvent("message", eventData)
			c.Writer.Flush()
		}
	}

	err = cozeConv.SendMessageStreamWithCallback(conversation.CozeConversationID, conversation.UserId, historyMessageList, onMessage)
	if err != nil {
		onMessage("error", map[string]string{"message": err.Error()})
		return
	}

	// 发送结束事件
	onMessage("end", map[string]string{"status": "completed"})
}

// 辅助函数：生成消息ID
func generateMessageId() string {
	return fmt.Sprintf("msg_%d", utils.GenerateSnowflakeId())
}

// 辅助函数：获取消息内容
func getMessageContent(msg interface{}) string {
	// 根据实际的消息结构提取内容
	if m, ok := msg.(map[string]interface{}); ok {
		if content, ok := m["content"].(string); ok {
			return content
		}
	}

	// 尝试从结构体中提取内容
	if msgStruct, ok := msg.(struct{ Content string }); ok {
		return msgStruct.Content
	}

	// 如果是文本形式，直接返回
	if content, ok := msg.(string); ok {
		return content
	}

	return "AI回复"
}

// 辅助函数：安全地将对象转换为 JSON
func mustMarshalJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return `{"error": "序列化失败"}`
	}
	return string(data)
}

func SendMessageWorkFlow(c *gin.Context) {
	// conversationIdStr := c.DefaultQuery("conversation_id", "0")
	// if conversationIdStr == "" {
	// 	conversationIdStr = "0"
	// }

	// conversationId, err := strconv.ParseUint(conversationIdStr, 10, 32)
	// if err != nil {
	// 	// 如果解析失败，默认为0
	// 	conversationId = 0
	// }

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数格式错误: "+err.Error())
		return
	}

	// userIdStr := c.DefaultQuery("user_id", "0")
	// if userIdStr == "" {
	// 	userIdStr = "0"
	// }

	// userID, err := strconv.ParseUint(userIdStr, 10, 32)
	// if err != nil {
	// 	// 如果解析失败，默认为0
	// 	userID = 0
	// }

	cozeConv, err := coze.New()
	if err != nil {
		utils.BadRequest(c, "初始化Coze对话失败: "+err.Error())
		return
	}

	resp, err := cozeConv.RunWorkflow(req.Content)
	if err != nil {
		utils.BadRequest(c, "工作流运行失败: "+err.Error())
		return
	}

	utils.Success(c,resp)
}

func SendMessageWorkFlowStream(c *gin.Context) {
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数格式错误: "+err.Error())
		return
	}

	cozeConv, err := coze.New()
	if err != nil {
		utils.BadRequest(c, "初始化Coze对话失败: "+err.Error())
		return
	}

	clientGone := c.Request.Context().Done()

	var aiMessageContent strings.Builder
	var aiMessageId string
	var aiMessageTokens int


	// 定义流式回调函数
	onMessage := func(eventType string, data interface{}) {
		select {
		case <-clientGone:
			return // 客户端已断开连接
		default:
			eventData := map[string]interface{}{
				"type": eventType,
				"data": data,
			}

			// 处理消息增量更新
			if eventType == "message_delta" {
				if msgData, ok := data.(map[string]interface{}); ok {
					if content, ok := msgData["content"].(string); ok {
						aiMessageContent.WriteString(content)
					}
				}
			}

			// 处理对话完成
			if eventType == "chat_completed" {
				if msgData, ok := data.(map[string]interface{}); ok {
					if chatId, ok := msgData["chat_id"].(string); ok {
						aiMessageId = chatId
					}

					if usage, ok := msgData["usage"].(map[string]interface{}); ok {
						for _, token := range usage {
							aiMessageTokens += token.(int)
						}
					}
				}

				// 保存AI回复消息
				if aiMessageContent.Len() > 0 {
					aiMessage := &models.Message{
						CozeMessageId:  aiMessageId,
						ConversationId: 10,
						ModelId:        1,
						Role:           "assistant",
						Content:        aiMessageContent.String(),
						Tokens:         aiMessageTokens,
					}

					if err := messageService.CreateMessage(aiMessage); err != nil {
						// 错误处理，但不中断流式响应
						fmt.Printf("保存AI回复失败: %v\n", err)
					}
				}
			}

			c.SSEvent("message", eventData)
			c.Writer.Flush()
		}
	}

	err = cozeConv.RunWorkflowStream(req.Content, onMessage)
	if err != nil {
		onMessage("workflow_error", map[string]string{"message": err.Error()})
		return
	}

	// 发送结束事件
	onMessage("end", map[string]string{"status": "completed"})	
}
