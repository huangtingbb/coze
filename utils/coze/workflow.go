package coze

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/coze-dev/coze-go"
)


func (workflow *Client) RunWorkflow(message string) (*coze.RunWorkflowsResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	workflow_param := map[string]interface{}{
		"ip": 	 "180.172.177.97",
		"token": "",
		"input": message,
	}
	workflowReq := &coze.RunWorkflowsReq{
		WorkflowID:  workflow.Config.WorkFlowID,
		Parameters: workflow_param,
		IsAsync: false,
	}

	resp, err := workflow.Api.Workflows.Runs.Create(ctx, workflowReq)
	if err != nil {
		return nil, fmt.Errorf("发送消息失败: %v", err)
	}

	return resp, nil
}

func (workflow *Client) RunWorkflowStream(message string, onMessage func(eventType string, data interface{})) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	workflow_param := map[string]interface{}{
		"ip": 	 "180.172.177.97",
		"token": "",
		"input": message,
	}
	workflowReq := &coze.RunWorkflowsReq{
		WorkflowID:  workflow.Config.WorkFlowID,
		Parameters: workflow_param,
		IsAsync: false,
	}

	resp, err := workflow.Api.Workflows.Runs.Stream(ctx, workflowReq)
	if err != nil {
		return fmt.Errorf("发送消息失败: %v", err)
	}

	handleWorkflowStream(resp, onMessage)

	return  nil
}

func handleWorkflowStream(resp coze.Stream[coze.WorkflowEvent], onMessage func(eventType string, data interface{})){
	defer resp.Close()
	for {
		event, err := resp.Recv()
		if errors.Is(err, io.EOF) {
			// 流式结束
			onMessage("workflow_end", map[string]string{
				"status": "completed",
				"log_id": resp.Response().LogID(),
			})
			break
		}
		if err != nil {
			fmt.Printf("工作流接收事件失败: %v\n", err)
			return
		}

		switch event.Event {
		case coze.WorkflowEventTypeMessage:
			// 流式增量
			onMessage("message_delta", map[string]string{
				"status": "delta",
				"log_id": resp.Response().LogID(),
				"content": event.Message.Content,
			})
			fmt.Printf("步骤开始: %s\n", event.Message.Content)
		case coze.WorkflowEventTypeError:
			onMessage("workflow_error", map[string]string{
				"status": "error",
				"log_id": resp.Response().LogID(),
				"content": event.Error.ErrorMessage,
			})
			fmt.Printf("步骤结束\n",)
		case coze.WorkflowEventTypeDone:
			onMessage("workflow_complated", map[string]string{
				"status": "completed",
				"log_id": resp.Response().LogID(),
				"content": "",
			})
			fmt.Println("流式响应结束")
			return
		default:
			fmt.Printf("未知事件: %v\n", event)
		}
	}
}
