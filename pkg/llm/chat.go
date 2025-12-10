package llm

import (
	"github.com/pandodao/tokenizer-go"
	openai "github.com/sashabaranov/go-openai"

	"github.com/eryajf/chatgpt-dingtalk/public"
)

// ChatWithContext 对话接口
func (c *Client) ChatWithContext(question string) (string, error) {
	if tokenizer.MustCalToken(question) > c.maxQuestionLen {
		return "", ErrOverMaxQuestionLength
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
	}

	resp, err := c.client.CreateChatCompletion(c.ctx, req)
	if err != nil {
		return "", err
	}

	answer := resp.Choices[0].Message.Content

	// 保存对话上下文
	c.ChatContext.old = append(c.ChatContext.old,
		conversation{Role: c.ChatContext.humanRole, Prompt: question},
		conversation{Role: c.ChatContext.aiRole, Prompt: answer},
	)
	c.ChatContext.seqTimes++

	return answer, nil
}
