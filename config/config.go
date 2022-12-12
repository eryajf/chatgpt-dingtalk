package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/eryajf/chatgpt-dingtalk/public/logger"
)

// Configuration 项目配置
type Configuration struct {
	// gtp apikey
	ApiKey string `json:"api_key"`
	// 会话超时时间
	SessionTimeout time.Duration `json:"session_timeout"`
	// GPT请求最大字符数
	MaxTokens uint `json:"max_tokens"`
	// GPT模型
	Model string `json:"model"`
	// 热度
	Temperature float64 `json:"temperature"`
	// 自定义清空会话口令
	SessionClearToken string `json:"session_clear_token"`
}

var config *Configuration
var once sync.Once

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	once.Do(func() {
		// 从文件中读取
		config = &Configuration{
			SessionTimeout:    60,
			MaxTokens:         512,
			Model:             "text-davinci-003",
			Temperature:       0.9,
			SessionClearToken: "下一个问题",
		}
		f, err := os.Open("config.json")
		if err != nil {
			logger.Danger("open config err: %v", err)
			return
		}
		defer f.Close()
		encoder := json.NewDecoder(f)
		err = encoder.Decode(config)
		if err != nil {
			logger.Warning("decode config err: %v", err)
			return
		}

		// 如果环境变量有配置，读取环境变量
		// 有环境变量使用环境变量
		ApiKey := os.Getenv("APIKEY")
		SessionTimeout := os.Getenv("SESSION_TIMEOUT")
		Model := os.Getenv("MODEL")
		MaxTokens := os.Getenv("MAX_TOKENS")
		Temperature := os.Getenv("TEMPREATURE")
		SessionClearToken := os.Getenv("SESSION_CLEAR_TOKEN")
		if ApiKey != "" {
			config.ApiKey = ApiKey
		}
		if SessionTimeout != "" {
			duration, err := time.ParseDuration(SessionTimeout)
			if err != nil {
				logger.Danger(fmt.Sprintf("config session timeout err: %v ,get is %v", err, SessionTimeout))
				return
			}
			config.SessionTimeout = duration
		}
		if Model != "" {
			config.Model = Model
		}
		if MaxTokens != "" {
			max, err := strconv.Atoi(MaxTokens)
			if err != nil {
				logger.Danger(fmt.Sprintf("config MaxTokens err: %v ,get is %v", err, MaxTokens))
				return
			}
			config.MaxTokens = uint(max)
		}
		if Temperature != "" {
			temp, err := strconv.ParseFloat(Temperature, 64)
			if err != nil {
				logger.Danger(fmt.Sprintf("config Temperature err: %v ,get is %v", err, Temperature))
				return
			}
			config.Temperature = temp
		}
		if SessionClearToken != "" {
			config.SessionClearToken = SessionClearToken
		}
	})
	if config.ApiKey == "" {
		logger.Danger("config err: api key required")
	}
	return config
}
