package chatgpt

import (
	"time"

	"github.com/avast/retry-go"
	"github.com/eryajf/chatgpt-dingtalk/public"
	"github.com/eryajf/chatgpt-dingtalk/public/logger"
)

func SingleQa(question, userId string) (answer string, err error) {
	chat := New(userId)
	defer chat.Close()
	// 定义一个重试策略
	retryStrategy := []retry.Option{
		retry.Delay(100 * time.Millisecond),
		retry.Attempts(3),
		retry.LastErrorOnly(true),
	}
	// 使用重试策略进行重试
	err = retry.Do(
		func() error {
			answer, err = chat.ChatWithContext(question)
			if err != nil {
				return err
			}
			return nil
		},
		retryStrategy...)
	return
}

func ContextQa(question, userId string) (chat *ChatGPT, answer string, err error) {
	chat = New(userId)
	if public.UserService.GetUserSessionContext(userId) != "" {
		err := chat.ChatContext.LoadConversation(userId)
		if err != nil {
			logger.Warning("load station failed: %v\n", err)
		}
	}
	retryStrategy := []retry.Option{
		retry.Delay(100 * time.Millisecond),
		retry.Attempts(3),
		retry.LastErrorOnly(true)}
	// 使用重试策略进行重试
	err = retry.Do(
		func() error {
			answer, err = chat.ChatWithContext(question)
			if err != nil {
				return err
			}
			return nil
		},
		retryStrategy...)
	return
}
