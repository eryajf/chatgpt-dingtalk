package llm

import (
	"context"
	"time"

	"github.com/avast/retry-go"

	"github.com/eryajf/chatgpt-dingtalk/pkg/logger"
	"github.com/eryajf/chatgpt-dingtalk/public"
)

// SingleQa 单聊
func SingleQa(question, userId string) (answer string, err error) {
	client := NewClient(userId)
	defer client.Close()

	retryStrategy := []retry.Option{
		retry.Delay(100 * time.Millisecond),
		retry.Attempts(3),
		retry.LastErrorOnly(true),
	}

	err = retry.Do(
		func() error {
			answer, err = client.ChatWithContext(question)
			if err != nil {
				return err
			}
			return nil
		},
		retryStrategy...)
	return
}

// ContextQa 串聊
func ContextQa(question, userId string) (client *Client, answer string, err error) {
	client = NewClient(userId)
	if public.UserService.GetUserSessionContext(userId) != "" {
		err := client.ChatContext.LoadConversation(userId)
		if err != nil {
			logger.Warning("load station failed: %v\n", err)
		}
	}
	retryStrategy := []retry.Option{
		retry.Delay(100 * time.Millisecond),
		retry.Attempts(3),
		retry.LastErrorOnly(true)}

	err = retry.Do(
		func() error {
			answer, err = client.ChatWithContext(question)
			if err != nil {
				return err
			}
			return nil
		},
		retryStrategy...)
	return
}

// ImageQa 生成图片
func ImageQa(ctx context.Context, question, userId string) (answer string, err error) {
	client := NewClient(userId)
	defer client.Close()

	retryStrategy := []retry.Option{
		retry.Delay(100 * time.Millisecond),
		retry.Attempts(3),
		retry.LastErrorOnly(true),
	}

	err = retry.Do(
		func() error {
			answer, err = client.GenerateImage(ctx, question)
			if err != nil {
				return err
			}
			return nil
		},
		retryStrategy...)
	return
}
