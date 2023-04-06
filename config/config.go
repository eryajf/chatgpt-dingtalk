package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eryajf/chatgpt-dingtalk/pkg/logger"
	"gopkg.in/yaml.v2"
)

// Configuration 项目配置
type Configuration struct {
	// 日志级别，info或者debug
	LogLevel string `yaml:"log_level"`
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
	// 限定对话类型 0：不限 1：单聊 2：群聊
	ChatType string `yaml:"chat_type"`
	// 哪些群组可以进行对话
	AllowGroups []string `yaml:"allow_groups"`
	// 哪些用户可以进行对话
	AllowUsers []string `yaml:"allow_users"`
	// 指定哪些人为此系统的管理员，必须指定，否则所有人都是
	AdminUsers []string `yaml:"admin_users"`
	// 钉钉机器人在应用信息中的AppSecret，为了校验回调的请求是否合法，如果你的服务对接给多个机器人，这里可以配置多个机器人的secret
	AppSecrets []string `yaml:"app_secrets"`
	// 自定义帮助信息
	Help string `yaml:"help"`
	// AzureOpenAI 配置
	AzureOn             bool   `yaml:"azure_on"`
	AzureApiVersion     string `yaml:"azure_api_version"`
	AzureResourceName   string `yaml:"azure_resource_name"`
	AzureDeploymentName string `yaml:"azure_deployment_name"`
	AzureOpenAIToken    string `yaml:"azure_openai_token"`
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
		logLevel := os.Getenv("LOG_LEVEL")
		if logLevel != "" {
			config.LogLevel = logLevel
		}
		apiKey := os.Getenv("APIKEY")
		if apiKey != "" {
			config.ApiKey = apiKey
		}
		baseURL := os.Getenv("BASE_URL")
		if baseURL != "" {
			config.BaseURL = baseURL
		}
		model := os.Getenv("MODEL")
		if model != "" {
			config.Model = model
		}
		sessionTimeout := os.Getenv("SESSION_TIMEOUT")
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
		defaultMode := os.Getenv("DEFAULT_MODE")
		if defaultMode != "" {
			config.DefaultMode = defaultMode
		}
		httpProxy := os.Getenv("HTTP_PROXY")
		if httpProxy != "" {
			config.HttpProxy = httpProxy
		}
		maxRequest := os.Getenv("MAX_REQUEST")
		if maxRequest != "" {
			newMR, _ := strconv.Atoi(maxRequest)
			config.MaxRequest = newMR
		}
		port := os.Getenv("PORT")
		if port != "" {
			config.Port = port
		}
		serviceURL := os.Getenv("SERVICE_URL")
		if serviceURL != "" {
			config.ServiceURL = serviceURL
		}
		chatType := os.Getenv("CHAT_TYPE")
		if chatType != "" {
			config.ChatType = chatType
		}
		allowGroup := os.Getenv("ALLOW_GROUPS")
		if allowGroup != "" {
			config.AllowGroups = strings.Split(allowGroup, ",")
		}
		allowUsers := os.Getenv("ALLOW_USERS")
		if allowUsers != "" {
			config.AllowUsers = strings.Split(allowUsers, ",")
		}
		adminUsers := os.Getenv("ADMIN_USERS")
		if adminUsers != "" {
			config.AdminUsers = strings.Split(adminUsers, ",")
		}
		appSecrets := os.Getenv("APP_SECRETS")
		if appSecrets != "" {
			config.AppSecrets = strings.Split(appSecrets, ",")
		}
		help := os.Getenv("HELP")
		if help != "" {
			config.Help = help
		}
		azureOn := os.Getenv("AZURE_ON")
		if azureOn != "" {
			config.AzureOn = azureOn == "true"
		}
		azureApiVersion := os.Getenv("AZURE_API_VERSION")
		if azureApiVersion != "" {
			config.AzureApiVersion = azureApiVersion
		}
		azureResourceName := os.Getenv("AZURE_RESOURCE_NAME")
		if azureResourceName != "" {
			config.AzureResourceName = azureResourceName
		}
		azureDeploymentName := os.Getenv("AZURE_DEPLOYMENT_NAME")
		if azureDeploymentName != "" {
			config.AzureDeploymentName = azureDeploymentName
		}
		azureOpenaiToken := os.Getenv("AZURE_OPENAI_TOKEN")
		if azureOpenaiToken != "" {
			config.AzureOpenAIToken = azureOpenaiToken
		}

	})

	// 一些默认值
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}
	if config.Model == "" {
		config.Model = "gpt-3.5-turbo"
	}
	if config.DefaultMode == "" {
		config.DefaultMode = "单聊"
	}
	if config.Port == "" {
		config.Port = "8090"
	}
	if config.ChatType == "" {
		config.ChatType = "0"
	}
	if config.ApiKey == "" {
		logger.Fatal("config err: api key required")
	}
	if config.ServiceURL == "" {
		logger.Fatal("config err: service url required")
	}
	return config
}
