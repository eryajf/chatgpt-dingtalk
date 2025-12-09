package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/eryajf/chatgpt-dingtalk/pkg/logger"
)

type Credential struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

// Configuration 项目配置
type Configuration struct {
	// 日志级别，info或者debug
	LogLevel string `yaml:"log_level"`
	// gpt apikey
	ApiKey string `yaml:"api_key"`
	// 运行模式
	RunMode string `yaml:"run_mode"`
	// 请求的 URL 地址
	BaseURL string `yaml:"base_url"`
	// 使用模型
	Model string `yaml:"model"`
	// 使用绘画模型
	ImageModel string `yaml:"image_model"`
	// 会话超时时间
	SessionTimeout time.Duration `yaml:"session_timeout"`
	// 最大问题长度
	MaxQuestionLen int `yaml:"max_question_len"`
	// 最大答案长度
	MaxAnswerLen int `yaml:"max_answer_len"`
	// 最大文本 = 问题 + 回答, 接口限制
	MaxText int `yaml:"max_text"`
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
	// 哪些outgoing群组可以进行对话
	AllowOutgoingGroups []string `yaml:"allow_outgoing_groups"`
	// 哪些用户可以进行对话
	AllowUsers []string `yaml:"allow_users"`
	// 哪些用户不可以进行对话
	DenyUsers []string `yaml:"deny_users"`
	// 哪些Vip用户可以进行无限对话
	VipUsers []string `yaml:"vip_users"`
	// 指定哪些人为此系统的管理员，必须指定，否则所有人都是
	AdminUsers []string `yaml:"admin_users"`
	// 钉钉机器人在应用信息中的AppSecret，为了校验回调的请求是否合法，如果你的服务对接给多个机器人，这里可以配置多个机器人的secret
	AppSecrets []string `yaml:"app_secrets"`
	// 敏感词，提问时触发，则不允许提问，回答的内容中触发，则以 🚫 代替
	SensitiveWords []string `yaml:"sensitive_words"`
	// 自定义帮助信息
	Help string `yaml:"help"`
	// AzureOpenAI 配置
	AzureOn             bool   `yaml:"azure_on"`
	AzureApiVersion     string `yaml:"azure_api_version"`
	AzureResourceName   string `yaml:"azure_resource_name"`
	AzureDeploymentName string `yaml:"azure_deployment_name"`
	AzureOpenAIToken    string `yaml:"azure_openai_token"`
	// 钉钉应用鉴权凭据
	Credentials []Credential `yaml:"credentials"`
	// 是否启用流式输出
	StreamMode bool `yaml:"stream_mode"`
	// 钉钉卡片模板ID(用于流式输出)
	CardTemplateID string `yaml:"card_template_id"`
}

var (
	config *Configuration
	once   sync.Once
)

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	once.Do(func() {
		// 从文件中读取
		config = &Configuration{}
		data, err := os.ReadFile("config.yml")
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
		runMode := os.Getenv("RUN_MODE")
		if runMode != "" {
			config.RunMode = runMode
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
		maxQuestionLen := os.Getenv("MAX_QUESTION_LEN")
		if maxQuestionLen != "" {
			newLen, _ := strconv.Atoi(maxQuestionLen)
			config.MaxQuestionLen = newLen
		}
		maxAnswerLen := os.Getenv("MAX_ANSWER_LEN")
		if maxAnswerLen != "" {
			newLen, _ := strconv.Atoi(maxAnswerLen)
			config.MaxAnswerLen = newLen
		}
		maxText := os.Getenv("MAX_TEXT")
		if maxText != "" {
			newLen, _ := strconv.Atoi(maxText)
			config.MaxText = newLen
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
		allowGroups := os.Getenv("ALLOW_GROUPS")
		if allowGroups != "" {
			config.AllowGroups = strings.Split(allowGroups, ",")
		}
		allowOutgoingGroups := os.Getenv("ALLOW_OUTGOING_GROUPS")
		if allowOutgoingGroups != "" {
			config.AllowOutgoingGroups = strings.Split(allowOutgoingGroups, ",")
		}
		allowUsers := os.Getenv("ALLOW_USERS")
		if allowUsers != "" {
			config.AllowUsers = strings.Split(allowUsers, ",")
		}
		denyUsers := os.Getenv("DENY_USERS")
		if denyUsers != "" {
			config.DenyUsers = strings.Split(denyUsers, ",")
		}
		vipUsers := os.Getenv("VIP_USERS")
		if vipUsers != "" {
			config.VipUsers = strings.Split(vipUsers, ",")
		}
		adminUsers := os.Getenv("ADMIN_USERS")
		if adminUsers != "" {
			config.AdminUsers = strings.Split(adminUsers, ",")
		}
		appSecrets := os.Getenv("APP_SECRETS")
		if appSecrets != "" {
			config.AppSecrets = strings.Split(appSecrets, ",")
		}
		sensitiveWords := os.Getenv("SENSITIVE_WORDS")
		if sensitiveWords != "" {
			config.SensitiveWords = strings.Split(sensitiveWords, ",")
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

		credentials := os.Getenv("DINGTALK_CREDENTIALS")
		if credentials != "" {
			config.Credentials = []Credential{}
			for _, idSecret := range strings.Split(credentials, ",") {
				items := strings.SplitN(idSecret, ":", 2)
				if len(items) == 2 {
					config.Credentials = append(config.Credentials, Credential{ClientID: items[0], ClientSecret: items[1]})
				}
			}
		}
	})

	// 一些默认值
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}
	if config.RunMode == "" {
		config.RunMode = "http"
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
	if !config.AzureOn {
		if config.ApiKey == "" {
			panic("config err: api key required")
		}
	}
	if config.MaxQuestionLen == 0 {
		config.MaxQuestionLen = 4096
	}
	if config.MaxAnswerLen == 0 {
		config.MaxAnswerLen = 4096
	}
	if config.MaxText == 0 {
		config.MaxText = 4096
	}
	return config
}
