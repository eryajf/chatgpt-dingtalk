package llm

import (
	"context"

	"github.com/eryajf/chatgpt-dingtalk/public"
)

// SingleQa 单聊
func SingleQa(question, userId string) (string, error) {
	client := NewClient(userId)
	defer client.Close()

	return client.ChatWithContext(question)
}

// ContextQa 串聊
func ContextQa(question, userId string) (*Client, string, error) {
	client := NewClient(userId)
	if public.UserService.GetUserSessionContext(userId) != "" {
		_ = client.ChatContext.LoadConversation(userId)
	}

	answer, err := client.ChatWithContext(question)
	return client, answer, err
}

// ImageQa 生成图片
func ImageQa(ctx context.Context, question, userId string) (string, error) {
	client := NewClient(userId)
	defer client.Close()

	return client.GenerateImage(ctx, question)
}
