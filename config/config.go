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
	// 默认对话模式
	DefaultMode string `json:"default_mode"`
}

var config *Configuration
var once sync.Once

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	once.Do(func() {
		// 从文件中读取
		config = &Configuration{}
		f, err := os.Open("config.json")
		if err != nil {
			logger.Danger(fmt.Errorf("open config err: %+v", err))
			return
		}
		defer f.Close()
		encoder := json.NewDecoder(f)
		err = encoder.Decode(config)
		if err != nil {
			logger.Warning(fmt.Errorf("decode config err: %v", err))
			return
		}
		// 如果环境变量有配置，读取环境变量
		ApiKey := os.Getenv("APIKEY")
		SessionTimeout := os.Getenv("SESSION_TIMEOUT")
		defaultMode := os.Getenv("DEFAULT_MODE")
		// Model := os.Getenv("MODEL")
		// MaxTokens := os.Getenv("MAX_TOKENS")
		// Temperature := os.Getenv("TEMPREATURE")
		// SessionClearToken := os.Getenv("SESSION_CLEAR_TOKEN")
		if ApiKey != "" {
			config.ApiKey = ApiKey
		}
		if SessionTimeout != "" {
			duration, err := strconv.ParseInt(SessionTimeout, 10, 64)
			if err != nil {
				logger.Danger(fmt.Sprintf("config session timeout err: %v ,get is %v", err, SessionTimeout))
				return
			}
			config.SessionTimeout = time.Duration(duration) * time.Second
		} else {
			config.SessionTimeout = time.Duration(config.SessionTimeout) * time.Second
		}
		if defaultMode != "" {
			config.DefaultMode = defaultMode
		}
	})
	if config.DefaultMode == "" {
		config.DefaultMode = "单聊"
	}
	if config.ApiKey == "" {
		logger.Danger("config err: api key required")
	}
	return config
}
