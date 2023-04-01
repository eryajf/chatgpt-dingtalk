package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/eryajf/chatgpt-dingtalk/pkg/logger"
	"gopkg.in/yaml.v2"
)

// Configuration 项目配置
type Configuration struct {
	// gtp apikey
	ApiKey string `yaml:"api_key"`
	// 请求的 URL 地址
	BaseURL string `yaml:"base_url"`
	// 使用模型
	Model string `yaml:"model"`
	// 会话超时时间
	SessionTimeout time.Duration `yaml:"session_timeout"`
	// 默认对话模式
	DefaultMode string `yaml:"default_mode"`
	// 代理地址
	HttpProxy string `yaml:"http_proxy"`
	// 用户单日最大请求次数
	MaxRequest int `yaml:"max_request"`
	// 指定服务启动端口，默认为 8090
	Port string `yaml:"port"`
	// 指定服务的地址，就是钉钉机器人配置的回调地址，比如: http://chat.eryajf.net
	ServiceURL string `yaml:"service_url"`
}

var config *Configuration
var once sync.Once

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	once.Do(func() {
		// 从文件中读取
		config = &Configuration{}
		data, err := ioutil.ReadFile("config.yml")
		if err != nil {
			log.Fatal(err)
		}
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			log.Fatal(err)
		}

		// 如果环境变量有配置，读取环境变量
		apiKey := os.Getenv("APIKEY")
		baseURL := os.Getenv("BASE_URL")
		model := os.Getenv("MODEL")
		sessionTimeout := os.Getenv("SESSION_TIMEOUT")
		defaultMode := os.Getenv("DEFAULT_MODE")
		httpProxy := os.Getenv("HTTP_PROXY")
		maxRequest := os.Getenv("MAX_REQUEST")
		port := os.Getenv("PORT")
		serviceURL := os.Getenv("SERVICE_URL")
		if apiKey != "" {
			config.ApiKey = apiKey
		}
		if baseURL != "" {
			config.BaseURL = baseURL
		}
		if sessionTimeout != "" {
			duration, err := strconv.ParseInt(sessionTimeout, 10, 64)
			if err != nil {
				logger.Fatal(fmt.Sprintf("config session timeout err: %v ,get is %v", err, sessionTimeout))
				return
			}
			config.SessionTimeout = time.Duration(duration) * time.Second
		} else {
			config.SessionTimeout = time.Duration(config.SessionTimeout) * time.Second
		}
		if defaultMode != "" {
			config.DefaultMode = defaultMode
		}
		if httpProxy != "" {
			config.HttpProxy = httpProxy
		}
		if model != "" {
			config.Model = model
		}
		if maxRequest != "" {
			newMR, _ := strconv.Atoi(maxRequest)
			config.MaxRequest = newMR
		}
		if port != "" {
			config.Port = port
		}
		if serviceURL != "" {
			config.ServiceURL = serviceURL
		}
	})
	if config.Model == "" {
		config.DefaultMode = "gpt-3.5-turbo"
	}
	if config.DefaultMode == "" {
		config.DefaultMode = "单聊"
	}
	if config.Port == "" {
		config.Port = "8090"
	}
	if config.ApiKey == "" {
		logger.Fatal("config err: api key required")
	}
	if config.ServiceURL == "" {
		logger.Fatal("config err: service url required")
	}
	return config
}
