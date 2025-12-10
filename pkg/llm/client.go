package llm

import (
	"context"
	"net/http"
	"net/url"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"github.com/eryajf/chatgpt-dingtalk/public"
)

type Client struct {
	client         *openai.Client
	ctx            context.Context
	userId         string
	maxQuestionLen int
	maxText        int
	maxAnswerLen   int
	timeOut        time.Duration
	doneChan       chan struct{}
	cancel         func()

	ChatContext *Context
}

func NewClient(userId string) *Client {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	timeOutChan := make(chan struct{}, 1)
	go func() {
		<-ctx.Done()
		timeOutChan <- struct{}{}
	}()

	config := openai.DefaultConfig(public.Config.ApiKey)

	// Azure配置
	if public.Config.AzureOn {
		config = openai.DefaultAzureConfig(
			public.Config.AzureOpenAIToken,
			"https://"+public.Config.AzureResourceName+".openai.azure.com",
		)
		config.APIVersion = public.Config.AzureApiVersion
		config.AzureModelMapperFunc = func(model string) string {
			return public.Config.AzureDeploymentName
		}
	} else {
		// HTTP客户端配置
		transport := &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		}

		if public.Config.HttpProxy != "" {
			proxyURL, _ := url.Parse(public.Config.HttpProxy)
			transport.Proxy = http.ProxyURL(proxyURL)
		}

		config.HTTPClient = &http.Client{Transport: transport}

		if public.Config.BaseURL != "" {
			config.BaseURL = public.Config.BaseURL + "/v1"
		}
	}

	return &Client{
		client:         openai.NewClientWithConfig(config),
		ctx:            ctx,
		userId:         userId,
		maxQuestionLen: public.Config.MaxQuestionLen,
		maxAnswerLen:   public.Config.MaxAnswerLen,
		maxText:        public.Config.MaxText,
		timeOut:        public.Config.SessionTimeout,
		doneChan:       timeOutChan,
		cancel:         cancel,
		ChatContext:    NewContext(),
	}
}

func (c *Client) Close() {
	c.cancel()
}

func (c *Client) GetDoneChan() chan struct{} {
	return c.doneChan
}

func (c *Client) SetMaxQuestionLen(maxQuestionLen int) int {
	if maxQuestionLen > c.maxText-c.maxAnswerLen {
		maxQuestionLen = c.maxText - c.maxAnswerLen
	}
	c.maxQuestionLen = maxQuestionLen
	return c.maxQuestionLen
}
