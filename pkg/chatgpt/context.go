package chatgpt

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"

	"golang.org/x/image/webp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"

	"os"
	"strings"
	"time"

	"github.com/pandodao/tokenizer-go"
	openai "github.com/sashabaranov/go-openai"

	"github.com/eryajf/chatgpt-dingtalk/pkg/dingbot"
	"github.com/eryajf/chatgpt-dingtalk/public"
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

// 通过 base64 编码字符串开头字符判断图像类型
func getImageTypeFromBase64(base64Str string) string {
	switch {
	case strings.HasPrefix(base64Str, "/9j/"):
		return "JPEG"
	case strings.HasPrefix(base64Str, "iVBOR"):
		return "PNG"
	case strings.HasPrefix(base64Str, "R0lG"):
		return "GIF"
	case strings.HasPrefix(base64Str, "UklG"):
		return "WebP"
	default:
		return "Unknown"
	}
}

func (c *ChatGPT) ChatWithContext(question string) (answer string, err error) {
	question = question + "."
	if tokenizer.MustCalToken(question) > c.maxQuestionLen {
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
	// 删除对话，直到prompt的长度满足条件
	for tokenizer.MustCalToken(prompt) > c.maxText {
		if len(c.ChatContext.old) > 1 { // 至少保留一条记录
			c.ChatContext.PollConversation() // 删除最旧的一条对话
			// 重新构建 prompt，计算长度
			promptTable = promptTable[1:] // 删除promptTable中对应的对话
			prompt = strings.Join(promptTable, "\n") + c.ChatContext.startSeq
		} else {
			break // 如果已经只剩一条记录，那么跳出循环
		}
	}
	//	if tokenizer.MustCalToken(prompt) > c.maxText-c.maxAnswerLen {
	//		return "", OverMaxTextLength
	//	}
	model := public.Config.Model
	userId := c.userId
	if public.Config.AzureOn {
		userId = ""
	}
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
func (c *ChatGPT) GenerateImage(ctx context.Context, prompt string) (string, error) {
	model := public.Config.Model
	imageModel := public.Config.ImageModel
	if isModelSupportedChatCompletions(model) {
		req := openai.ImageRequest{
			Prompt:         prompt,
			Model:          imageModel,
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

		// dall-e-3 返回的是 WebP 格式的图片，需要判断处理
		imgType := getImageTypeFromBase64(respBase64.Data[0].B64JSON)
		var imgData image.Image
		var imgErr error
		if imgType == "WebP" {
			imgData, imgErr = webp.Decode(r)
		} else {
			imgData, _, imgErr = image.Decode(r)
		}
		if imgErr != nil {
			return "", imgErr
		}

		imageName := time.Now().Format("20060102-150405") + ".png"
		clientId, _ := ctx.Value(public.DingTalkClientIdKeyName).(string)
		client := public.DingTalkClientManager.GetClientByOAuthClientID(clientId)
		mediaResult, uploadErr := &dingbot.MediaUploadResult{}, errors.New(fmt.Sprintf("unknown clientId: %s", clientId))
		if client != nil {
			mediaResult, uploadErr = client.UploadMedia(imgBytes, imageName, dingbot.MediaTypeImage, dingbot.MimeTypeImagePng)
		}

		err = os.MkdirAll("data/images", 0755)
		if err != nil {
			return "", err
		}
		file, err := os.Create("data/images/" + imageName)
		if err != nil {
			return "", err
		}
		defer file.Close()

		if err := png.Encode(file, imgData); err != nil {
			return "", err
		}
		if uploadErr == nil {
			return mediaResult.MediaID, nil
		} else {
			return public.Config.ServiceURL + "/images/" + imageName, nil
		}
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
