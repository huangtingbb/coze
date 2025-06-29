package coze

import (
	"coze-agent-platform/config"
	"fmt"
	"context"

	"github.com/coze-dev/coze-go"
)

type Conversation struct {
	Config *config.CozeConfig
	Cli *coze.CozeAPI
}

func NewConversation() (*Conversation, error)  {
	cozeConv := &Conversation{
		Config: config.GetCozeConfig(),
	}
	token, err := GetToken()
	if err != nil {
		return nil, fmt.Errorf("获取Coze Token失败: %v", err)
	}

	cozeCli := coze.NewCozeAPI(coze.NewTokenAuth(token),coze.WithBaseURL(cozeConv.Config.APIURL))
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

func (conversation *Conversation) SendMessage(conversationID string, message string) (*coze.CreateMessageResp, error){
	messageReq := &coze.CreateMessageReq{}
	messageReq.ConversationID = conversationID
	messageReq.Role = coze.MessageRoleUser
	messageReq.SetObjectContext([]*coze.MessageObjectString{
		coze.NewTextMessageObject(message),
	})
	ctx := context.Background()
	msgs, err := conversation.Cli.Conversations.Messages.Create(ctx, messageReq)
	if err != nil {
		return nil,fmt.Errorf("发送消息失败: %v", err)
	}

	return msgs,nil
}

