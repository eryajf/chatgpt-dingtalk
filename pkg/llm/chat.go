package llm

import (
	"fmt"
	"strings"

	"github.com/pandodao/tokenizer-go"
	openai "github.com/sashabaranov/go-openai"

	"github.com/eryajf/chatgpt-dingtalk/public"
)

func (c *Client) ChatWithContext(question string) (answer string, err error) {
	question = question + "."
	if tokenizer.MustCalToken(question) > c.maxQuestionLen {
		return "", ErrOverMaxQuestionLength
	}
	if c.ChatContext.seqTimes >= c.ChatContext.maxSeqTimes {
		if c.ChatContext.maintainSeqTimes {
			c.ChatContext.PollConversation()
		} else {
			return "", ErrOverMaxSequenceTimes
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
	fmt.Println("Using model:", model)
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
		}
		resp, err := c.client.CreateChatCompletion(c.ctx, req)
		if err != nil {
			return "", err
		}
		resp.Choices[0].Message.Content = formatAnswer(resp.Choices[0].Message.Content)
		c.ChatContext.old = append(c.ChatContext.old, conversation{
			Role:   c.ChatContext.humanRole,
			Prompt: question,
		})
		c.ChatContext.old = append(c.ChatContext.old, conversation{
			Role:   c.ChatContext.aiRole,
			Prompt: resp.Choices[0].Message.Content,
		})
		c.ChatContext.seqTimes++
		return resp.Choices[0].Message.Content, nil
	} else {
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
			return "", err
		}
		resp.Choices[0].Text = formatAnswer(resp.Choices[0].Text)
		c.ChatContext.old = append(c.ChatContext.old, conversation{
			Role:   c.ChatContext.humanRole,
			Prompt: question,
		})
		c.ChatContext.old = append(c.ChatContext.old, conversation{
			Role:   c.ChatContext.aiRole,
			Prompt: resp.Choices[0].Text,
		})
		c.ChatContext.seqTimes++
		return resp.Choices[0].Text, nil
	}
}
