package gpt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/eryajf/chatgpt-dingtalk/config"
	"github.com/eryajf/chatgpt-dingtalk/public/logger"
	"github.com/go-resty/resty/v2"
)

const BASEURL = "https://api.openai.com/v1/"

// ChatGPTRequestBody 请求体
type ChatGPTRequestBody struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   uint    `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

// ChatGPTResponseBody 响应体
type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChoiceItem           `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}

type ChoiceItem struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	Logprobs     int    `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

// Completions gtp文本模型回复
//curl https://api.openai.com/v1/completions
//-H "Content-Type: application/json"
//-H "Authorization: Bearer your chatGPT key"
//-d '{"model": "text-davinci-003", "prompt": "give me good song", "temperature": 0, "max_tokens": 7}'
func Completions(msg string) (string, error) {
	cfg := config.LoadConfig()
	requestBody := ChatGPTRequestBody{
		Model:       cfg.Model,
		Prompt:      msg,
		MaxTokens:   cfg.MaxTokens,
		Temperature: cfg.Temperature,
	}

	client := resty.New().
		SetRetryCount(2).
		SetRetryWaitTime(1*time.Second).
		SetTimeout(cfg.SessionTimeout).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+cfg.ApiKey)

	rsp, err := client.R().SetBody(requestBody).Post(BASEURL + "completions")
	if err != nil {
		return "", fmt.Errorf("request openai failed, err : %v", err)
	}
	if rsp.StatusCode() != 200 {
		return "", fmt.Errorf("gtp api status code not equals 200, code is %d ,details:  %v ", rsp.StatusCode(), string(rsp.Body()))
	} else {
		logger.Info(fmt.Sprintf("response gtp json string : %v", string(rsp.Body())))
	}

	gptResponseBody := &ChatGPTResponseBody{}
	err = json.Unmarshal(rsp.Body(), gptResponseBody)
	if err != nil {
		return "", err
	}
	var reply string
	if len(gptResponseBody.Choices) > 0 {
		reply = gptResponseBody.Choices[0].Text
	}
	return reply, nil
}
