package chatgpt

import (
	"context"
	"net/http"
	"net/url"
	"time"

	gogpt "github.com/sashabaranov/go-gpt3"
)

type ChatGPT struct {
	client         *gogpt.Client
	ctx            context.Context
	userId         string
	maxQuestionLen int
	maxText        int
	maxAnswerLen   int
	timeOut        time.Duration // 超时时间, 0表示不超时
	doneChan       chan struct{}
	cancel         func()

	ChatContext *ChatContext
}

func New(apiKey, proxyUrl, userId string, timeOut time.Duration) *ChatGPT {
	var ctx context.Context
	var cancel func()
	if timeOut == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), timeOut)
	}
	timeOutChan := make(chan struct{}, 1)
	go func() {
		<-ctx.Done()
		timeOutChan <- struct{}{} // 发送超时信号，或是提示结束，用于聊天机器人场景，配合GetTimeOutChan() 使用
	}()

	config := gogpt.DefaultConfig(apiKey)
	if proxyUrl != "" {
		config.HTTPClient.Transport = &http.Transport{
			// 设置代理
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(proxyUrl)
			}}
	}
	return &ChatGPT{
		client:         gogpt.NewClientWithConfig(config),
		ctx:            ctx,
		userId:         userId,
		maxQuestionLen: 2048, // 最大问题长度
		maxAnswerLen:   2048, // 最大答案长度
		maxText:        4096, // 最大文本 = 问题 + 回答, 接口限制
		timeOut:        timeOut,
		doneChan:       timeOutChan,
		cancel: func() {
			cancel()
		},
		ChatContext: NewContext(),
	}
}
func (c *ChatGPT) Close() {
	c.cancel()
}

func (c *ChatGPT) GetDoneChan() chan struct{} {
	return c.doneChan
}

func (c *ChatGPT) SetMaxQuestionLen(maxQuestionLen int) int {
	if maxQuestionLen > c.maxText-c.maxAnswerLen {
		maxQuestionLen = c.maxText - c.maxAnswerLen
	}
	c.maxQuestionLen = maxQuestionLen
	return c.maxQuestionLen
}

// func (c *ChatGPT) Chat(question string) (answer string, err error) {
// 	question = question + "."
// 	if len(question) > c.maxQuestionLen {
// 		return "", OverMaxQuestionLength
// 	}
// 	if len(question)+c.maxAnswerLen > c.maxText {
// 		question = question[:c.maxText-c.maxAnswerLen]
// 	}
// 	req := gogpt.CompletionRequest{
// 		Model:            gogpt.GPT3TextDavinci003,
// 		MaxTokens:        c.maxAnswerLen,
// 		Prompt:           question,
// 		Temperature:      0.9,
// 		TopP:             1,
// 		N:                1,
// 		FrequencyPenalty: 0,
// 		PresencePenalty:  0.5,
// 		User:             c.userId,
// 		Stop:             []string{},
// 	}
// 	resp, err := c.client.CreateCompletion(c.ctx, req)
// 	if err != nil {
// 		return "", err
// 	}
// 	return formatAnswer(resp.Choices[0].Text), err
// }
