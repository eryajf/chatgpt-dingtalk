package llm

import (
	"errors"
	"io"

	"github.com/pandodao/tokenizer-go"
	openai "github.com/sashabaranov/go-openai"

	"github.com/eryajf/chatgpt-dingtalk/public"
)

// ChatWithContextStream 流式对话,返回一个channel用于接收流式内容
func (c *Client) ChatWithContextStream(question string) (<-chan string, error) {
	if tokenizer.MustCalToken(question) > c.maxQuestionLen {
		return nil, ErrOverMaxQuestionLength
	}

	// 构建消息列表
	messages := c.buildMessages(question)

	model := public.Config.Model
	userId := c.userId
	if public.Config.AzureOn {
		userId = ""
	}

	req := openai.ChatCompletionRequest{
		Model:       model,
		Messages:    messages,
		MaxTokens:   c.maxAnswerLen,
		Temperature: 0.6,
		User:        userId,
		Stream:      true,
	}

	contentCh := make(chan string, 10)

	go func() {
		defer close(contentCh)

		stream, err := c.client.CreateChatCompletionStream(c.ctx, req)
		if err != nil {
			contentCh <- err.Error()
			return
		}
		defer stream.Close()

		fullAnswer := ""
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				if fullAnswer == "" {
					contentCh <- err.Error()
				}
				return
			}

			if len(response.Choices) > 0 {
				delta := response.Choices[0].Delta.Content
				if delta != "" {
					fullAnswer += delta
					contentCh <- delta
				}
			}
		}

		// 保存对话上下文
		c.ChatContext.old = append(c.ChatContext.old,
			conversation{Role: c.ChatContext.humanRole, Prompt: question},
			conversation{Role: c.ChatContext.aiRole, Prompt: fullAnswer},
		)
		c.ChatContext.seqTimes++
	}()

	return contentCh, nil
}

// buildMessages 构建消息列表
func (c *Client) buildMessages(question string) []openai.ChatCompletionMessage {
	var messages []openai.ChatCompletionMessage

	// 添加历史对话
	for _, v := range c.ChatContext.old {
		role := "assistant"
		if v.Role == c.ChatContext.humanRole {
			role = "user"
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: v.Prompt,
		})
	}

	// 添加当前问题
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    "user",
		Content: question,
	})

	return messages
}

// SingleQaStream 单聊流式版本
func SingleQaStream(question, userId string) (<-chan string, func(), error) {
	client := NewClient(userId)

	contentCh := make(chan string, 10)
	done := make(chan struct{})

	go func() {
		defer close(contentCh)
		defer close(done)

		stream, err := client.ChatWithContextStream(question)
		if err != nil {
			contentCh <- err.Error()
			client.Close()
			return
		}

		for content := range stream {
			contentCh <- content
		}

		client.Close()
	}()

	cleanup := func() {
		<-done
	}

	return contentCh, cleanup, nil
}

// ContextQaStream 串聊流式版本
func ContextQaStream(question, userId string) (*Client, <-chan string, error) {
	client := NewClient(userId)
	if public.UserService.GetUserSessionContext(userId) != "" {
		_ = client.ChatContext.LoadConversation(userId)
	}

	stream, err := client.ChatWithContextStream(question)
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, stream, nil
}
