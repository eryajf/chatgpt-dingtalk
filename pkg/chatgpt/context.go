package chatgpt

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"image/png"
	"os"
	"strings"
	"time"

	"github.com/eryajf/chatgpt-dingtalk/public"
	openai "github.com/sashabaranov/go-openai"
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
		background:       "",
		maxSeqTimes:      1000,
		preset:           "",
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
	model := public.Config.Model
	if model == openai.GPT3Dot5Turbo0301 ||
		model == openai.GPT3Dot5Turbo ||
		model == openai.GPT4 || model == openai.GPT40314 ||
		model == openai.GPT432K || model == openai.GPT432K0314 {
		req := openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: prompt,
				},
			},
			MaxTokens:   3072,
			Temperature: 0.6,
			User:        c.userId,
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
func (c *ChatGPT) GenreateImage(prompt string) (string, error) {
	model := public.Config.Model
	if model == openai.GPT3Dot5Turbo0301 ||
		model == openai.GPT3Dot5Turbo ||
		model == openai.GPT4 || model == openai.GPT40314 ||
		model == openai.GPT432K || model == openai.GPT432K0314 {
		req := openai.ImageRequest{
			Prompt:         prompt,
			Size:           openai.CreateImageSize1024x1024,
			ResponseFormat: openai.CreateImageResponseFormatB64JSON,
			N:              1,
			User:           c.userId,
		}
		respBase64, err := c.client.CreateImage(c.ctx, req)
		if err != nil {
			return "", err
		}
		imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
		if err != nil {
			return "", err
		}

		r := bytes.NewReader(imgBytes)
		imgData, err := png.Decode(r)
		if err != nil {
			return "", err
		}

		imageName := time.Now().Format("20060102-150405") + ".png"
		err = os.MkdirAll("images", 0755)
		if err != nil {
			return "", err
		}
		file, err := os.Create("images/" + imageName)
		if err != nil {
			return "", err
		}
		defer file.Close()

		if err := png.Encode(file, imgData); err != nil {
			return "", err
		}

		return public.Config.ServiceURL + "/images/" + imageName, nil
	}
	return "", nil
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
