package chatgpt

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/eryajf/chatgpt-dingtalk/public"
	gogpt "github.com/sashabaranov/go-gpt3"
)

var (
	DefaultAiRole    = "AI"
	DefaultHumanRole = "Human"

	DefaultCharacter  = []string{"helpful", "creative", "clever", "friendly", "lovely", "talkative"}
	DefaultBackground = "The following is a conversation with AI assistant. The assistant is %s"
	DefaultPreset     = "\n%s: 你好，让我们开始愉快的谈话！\n%s: 我是 AI assistant ，请问你有什么问题？"
)

type (
	ChatContext struct {
		background  string // 对话背景
		preset      string // 预设对话
		maxSeqTimes int    // 最大对话次数
		aiRole      *role  // AI角色
		humanRole   *role  // 人类角色

		old        []conversation // 旧对话
		restartSeq string         // 重新开始对话的标识
		startSeq   string         // 开始对话的标识

		seqTimes int // 对话次数

		maintainSeqTimes bool // 是否维护对话次数 (自动移除旧对话)
	}

	ChatContextOption func(*ChatContext)

	conversation struct {
		Role   *role
		Prompt string
	}

	role struct {
		Name string
	}
)

func NewContext(options ...ChatContextOption) *ChatContext {
	ctx := &ChatContext{
		aiRole:           &role{Name: DefaultAiRole},
		humanRole:        &role{Name: DefaultHumanRole},
		background:       fmt.Sprintf(DefaultBackground, strings.Join(DefaultCharacter, ", ")+"."),
		maxSeqTimes:      1000,
		preset:           fmt.Sprintf(DefaultPreset, DefaultHumanRole, DefaultAiRole),
		old:              []conversation{},
		seqTimes:         0,
		restartSeq:       "\n" + DefaultHumanRole + ": ",
		startSeq:         "\n" + DefaultAiRole + ": ",
		maintainSeqTimes: false,
	}

	for _, option := range options {
		option(ctx)
	}
	return ctx
}

// PollConversation 移除最旧的一则对话
func (c *ChatContext) PollConversation() {
	c.old = c.old[1:]
	c.seqTimes--
}

// ResetConversation 重置对话
func (c *ChatContext) ResetConversation(userid string) {
	public.UserService.ClearUserSessionContext(userid)
}

// SaveConversation 保存对话
func (c *ChatContext) SaveConversation(userid string) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(c.old)
	if err != nil {
		return err
	}
	public.UserService.SetUserSessionContext(userid, buffer.String())
	return nil
}

// LoadConversation 加载对话
func (c *ChatContext) LoadConversation(userid string) error {
	dec := gob.NewDecoder(strings.NewReader(public.UserService.GetUserSessionContext(userid)))
	err := dec.Decode(&c.old)
	if err != nil {
		return err
	}
	c.seqTimes = len(c.old)
	return nil
}

func (c *ChatContext) SetHumanRole(role string) {
	c.humanRole.Name = role
	c.restartSeq = "\n" + c.humanRole.Name + ": "
}

func (c *ChatContext) SetAiRole(role string) {
	c.aiRole.Name = role
	c.startSeq = "\n" + c.aiRole.Name + ": "
}

func (c *ChatContext) SetMaxSeqTimes(times int) {
	c.maxSeqTimes = times
}

func (c *ChatContext) GetMaxSeqTimes() int {
	return c.maxSeqTimes
}

func (c *ChatContext) SetBackground(background string) {
	c.background = background
}

func (c *ChatContext) SetPreset(preset string) {
	c.preset = preset
}

func (c *ChatGPT) ChatWithContext(question string) (answer string, err error) {
	question = question + "."
	if len(question) > c.maxQuestionLen {
		return "", OverMaxQuestionLength
	}
	if c.ChatContext.seqTimes >= c.ChatContext.maxSeqTimes {
		if c.ChatContext.maintainSeqTimes {
			c.ChatContext.PollConversation()
		} else {
			return "", OverMaxSequenceTimes
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
	if len(prompt) > c.maxText-c.maxAnswerLen {
		return "", OverMaxTextLength
	}

	if public.Config.Model == gogpt.GPT3Dot5Turbo0301 || public.Config.Model == gogpt.GPT3Dot5Turbo {
		req := gogpt.ChatCompletionRequest{
			Model: public.Config.Model,
			Messages: []gogpt.ChatCompletionMessage{
				{
					Role:    "user",
					Content: prompt,
				},
			}}
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
		req := gogpt.CompletionRequest{
			Model:            public.Config.Model,
			MaxTokens:        c.maxAnswerLen,
			Prompt:           prompt,
			Temperature:      0.9,
			TopP:             1,
			N:                1,
			FrequencyPenalty: 0,
			PresencePenalty:  0.5,
			User:             c.userId,
			Stop:             []string{c.ChatContext.aiRole.Name + ":", c.ChatContext.humanRole.Name + ":"},
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

func WithMaxSeqTimes(times int) ChatContextOption {
	return func(c *ChatContext) {
		c.SetMaxSeqTimes(times)
	}
}

// WithOldConversation 从文件中加载对话
func WithOldConversation(userid string) ChatContextOption {
	return func(c *ChatContext) {
		_ = c.LoadConversation(userid)
	}
}

func WithMaintainSeqTimes(maintain bool) ChatContextOption {
	return func(c *ChatContext) {
		c.maintainSeqTimes = maintain
	}
}
