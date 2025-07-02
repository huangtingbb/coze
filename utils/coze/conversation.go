package coze

import (
	"context"
	"coze-agent-platform/config"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/coze-dev/coze-go"
)

type Conversation struct {
	Config *config.CozeConfig
	Cli    *coze.CozeAPI
}

func NewConversation() (*Conversation, error) {
	cozeConv := &Conversation{
		Config: config.GetCozeConfig(),
	}
	token, err := GetToken()
	if err != nil {
		return nil, fmt.Errorf("获取Coze Token失败: %v", err)
	}
	httpClient := &http.Client{
		Timeout: 120 * time.Second,
	}

	cozeCli := coze.NewCozeAPI(coze.NewTokenAuth(token), coze.WithBaseURL(cozeConv.Config.APIURL), coze.WithHttpClient(httpClient))
	cozeConv.Cli = &cozeCli
	return cozeConv, nil
}

func (conversation *Conversation) CreateConversation() (string, error) {
	botID := conversation.Config.BotID
	ctx := context.Background()
	resp, err := conversation.Cli.Conversations.Create(ctx, &coze.CreateConversationsReq{BotID: botID})
	if err != nil {
		return "", fmt.Errorf("创建对话失败: %v", err)
	}

	return resp.Conversation.ID, nil
}

func (conversation *Conversation) SendMessage(conversationID string, message string) (*coze.CreateMessageResp, error) {
	messageReq := &coze.CreateMessageReq{}
	messageReq.ConversationID = conversationID
	messageReq.Role = coze.MessageRoleUser
	messageReq.SetObjectContext([]*coze.MessageObjectString{
		coze.NewTextMessageObject(message),
	})
	ctx := context.Background()
	msgs, err := conversation.Cli.Conversations.Messages.Create(ctx, messageReq)
	if err != nil {
		return nil, fmt.Errorf("发送消息失败: %v", err)
	}

	return msgs, nil
}

func (conversation *Conversation) SendMessageStream(conversationID string, message string) (*coze.CreateMessageResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	req := &coze.CreateChatsReq{
		BotID:  conversation.Config.BotID,
		UserID: "157",
		Messages: []*coze.Message{
			coze.BuildUserQuestionText(message, nil),
		},
	}
	fmt.Println(req)
	resp, err := conversation.Cli.Chat.Stream(ctx, req)

	if err != nil {
		return nil, fmt.Errorf("创建流式对话失败: %v", err)
	}
	defer resp.Close()

	for {
		event, err := resp.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("流式对话结束")
			break
		}

		fmt.Println(event)

		if err != nil {
			return nil, fmt.Errorf("流式对话失败: %v", err)
		}

		if event.Event == coze.ChatEventConversationMessageDelta {
			fmt.Println(event.Message.Content)
		} else if event.Event == coze.ChatEventConversationChatCompleted {
			fmt.Printf("流式对话完成,token usage:%d\n", event.Chat.Usage.TokenCount)
		} else if event.Event == coze.ChatEventConversationChatFailed {
			fmt.Printf("流式对话错误: code %d , msg %s\n", event.Chat.LastError.Code, event.Chat.LastError.Msg)
		} else {
			fmt.Printf("流式对话事件: %s\n", event.Event)
		}

		fmt.Printf("done, log :s% \n", resp.Response().LogID())
	}
	return nil, nil
}

// SendMessageStreamWithCallback 发送流式消息并通过回调函数处理事件
func (conversation *Conversation) SendMessageStreamWithCallback(conversationID string, message string, onMessage func(eventType string, data interface{})) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	req := &coze.CreateChatsReq{
		BotID:  conversation.Config.BotID,
		UserID: "157",
		Messages: []*coze.Message{
			coze.BuildUserQuestionText(message, nil),
		},
	}

	resp, err := conversation.Cli.Chat.Stream(ctx, req)
	if err != nil {
		return fmt.Errorf("创建流式对话失败: %v", err)
	}
	defer resp.Close()

	for {
		event, err := resp.Recv()
		if errors.Is(err, io.EOF) {
			// 流式对话结束
			onMessage("conversation_end", map[string]string{
				"status": "completed",
				"log_id": resp.Response().LogID(),
			})
			break
		}

		if err != nil {
			return fmt.Errorf("流式对话失败: %v", err)
		}

		// 根据不同的事件类型调用回调函数
		switch event.Event {
		case coze.ChatEventConversationMessageDelta:
			// 消息增量更新
			if event.Message != nil {
				onMessage("message_delta", map[string]interface{}{
					"content": event.Message.Content,
					"role":    event.Message.Role,
					"type":    event.Message.Type,
				})
			}
		case coze.ChatEventConversationChatCompleted:
			// 对话完成
			onMessage("chat_completed", map[string]interface{}{
				"usage": map[string]interface{}{
					"token_count":  event.Chat.Usage.TokenCount,
					"output_count": event.Chat.Usage.OutputCount,
					"input_count":  event.Chat.Usage.InputCount,
				},
				"chat_id": event.Chat.ID,
			})
		case coze.ChatEventConversationChatFailed:
			// 对话失败
			onMessage("chat_failed", map[string]interface{}{
				"error_code": event.Chat.LastError.Code,
				"error_msg":  event.Chat.LastError.Msg,
			})
		case coze.ChatEventConversationChatRequiresAction:
			// 需要用户操作
			onMessage("requires_action", map[string]interface{}{
				"action": event.Chat.RequiredAction,
			})
		default:
			// 其他事件
			onMessage("other_event", map[string]interface{}{
				"event":    string(event.Event),
				"raw_data": event,
			})
		}
	}

	return nil
}
