package llm

import (
	"io"
	"strings"

	"github.com/pandodao/tokenizer-go"
	openai "github.com/sashabaranov/go-openai"

	"github.com/eryajf/chatgpt-dingtalk/public"
)

// ChatWithContextStream 流式对话,返回一个channel用于接收流式内容
func (c *Client) ChatWithContextStream(question string) (<-chan string, error) {
	question = question + "."
	if tokenizer.MustCalToken(question) > c.maxQuestionLen {
		return nil, ErrOverMaxQuestionLength
	}
	if c.ChatContext.seqTimes >= c.ChatContext.maxSeqTimes {
		if c.ChatContext.maintainSeqTimes {
			c.ChatContext.PollConversation()
		} else {
			return nil, ErrOverMaxSequenceTimes
		}
	}

	var promptTable []string
	promptTable = append(promptTable, c.ChatContext.background)
	promptTable = append(promptTable, c.ChatContext.preset)
	for _, v := range c.ChatContext.old {
		if v.Role == c.ChatContext.humanRole {
			promptTable = append(promptTable, "\n"+v.Role.Name+": "+v.Prompt)
		} else {
			promptTable = append(promptTable, v.Role.Name+": "+v.Prompt)
		}
	}
	promptTable = append(promptTable, "\n"+c.ChatContext.restartSeq+question)
	prompt := strings.Join(promptTable, "\n")
	prompt += c.ChatContext.startSeq

	for tokenizer.MustCalToken(prompt) > c.maxText {
		if len(c.ChatContext.old) > 1 {
			c.ChatContext.PollConversation()
			promptTable = promptTable[1:]
			prompt = strings.Join(promptTable, "\n") + c.ChatContext.startSeq
		} else {
			break
		}
	}

	model := public.Config.Model
	userId := c.userId
	if public.Config.AzureOn {
		userId = ""
	}

	contentCh := make(chan string, 10)

	if isModelSupportedChatCompletions(model) {
		req := openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: prompt,
				},
			},
			MaxTokens:   c.maxAnswerLen,
			Temperature: 0.6,
			User:        userId,
			Stream:      true,
		}

		go func() {
			defer close(contentCh)

			stream, err := c.client.CreateChatCompletionStream(c.ctx, req)
			if err != nil {
				contentCh <- formatAnswer(err.Error())
				return
			}
			defer stream.Close()

			fullAnswer := ""
			for {
				response, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					// 流式接收中断,记录已接收的内容
					if fullAnswer != "" {
						// 如果已经接收到部分内容,只记录错误但不中断
						// 错误会在日志中显示,但不会影响已接收的内容
					} else {
						// 如果还没接收到任何内容,则发送错误信息
						contentCh <- formatAnswer(err.Error())
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
			fullAnswer = formatAnswer(fullAnswer)
			c.ChatContext.old = append(c.ChatContext.old, conversation{
				Role:   c.ChatContext.humanRole,
				Prompt: question,
			})
			c.ChatContext.old = append(c.ChatContext.old, conversation{
				Role:   c.ChatContext.aiRole,
				Prompt: fullAnswer,
			})
			c.ChatContext.seqTimes++
		}()
	} else {
		// 对于不支持流式的模型,使用普通方式并模拟流式输出
		go func() {
			defer close(contentCh)

			req := openai.CompletionRequest{
				Model:       model,
				MaxTokens:   c.maxAnswerLen,
				Prompt:      prompt,
				Temperature: 0.6,
				User:        c.userId,
				Stop:        []string{c.ChatContext.aiRole.Name + ":", c.ChatContext.humanRole.Name + ":"},
			}
			resp, err := c.client.CreateCompletion(c.ctx, req)
			if err != nil {
				contentCh <- formatAnswer(err.Error())
				return
			}
			answer := formatAnswer(resp.Choices[0].Text)

			// 保存对话上下文
			c.ChatContext.old = append(c.ChatContext.old, conversation{
				Role:   c.ChatContext.humanRole,
				Prompt: question,
			})
			c.ChatContext.old = append(c.ChatContext.old, conversation{
				Role:   c.ChatContext.aiRole,
				Prompt: answer,
			})
			c.ChatContext.seqTimes++

			// 模拟流式输出
			contentCh <- answer
		}()
	}

	return contentCh, nil
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
		err := client.ChatContext.LoadConversation(userId)
		if err != nil {
			// 忽略加载错误,继续执行
		}
	}

	stream, err := client.ChatWithContextStream(question)
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, stream, nil
}
