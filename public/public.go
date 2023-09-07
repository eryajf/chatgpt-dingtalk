package public

import (
	"github.com/sashabaranov/go-openai"

	"github.com/eryajf/chatgpt-dingtalk/config"
	"github.com/eryajf/chatgpt-dingtalk/pkg/cache"
	"github.com/eryajf/chatgpt-dingtalk/pkg/db"
	"github.com/eryajf/chatgpt-dingtalk/pkg/dingbot"
)

var UserService cache.UserServiceInterface
var Config *config.Configuration
var Prompt *[]config.Prompt
var DingTalkClientManager dingbot.DingTalkClientManagerInterface

const DingTalkClientIdKeyName = "DingTalkClientId"

func InitSvc() {
	// 加载配置
	Config = config.LoadConfig()
	// 加载prompt
	Prompt = config.LoadPrompt()
	// 初始化缓存
	UserService = cache.NewUserService()
	// 初始化钉钉开放平台的客户端，用于访问上传图片等能力
	DingTalkClientManager = dingbot.NewDingTalkClientManager(Config)
	// 初始化数据库
	db.InitDB()
	// 暂时不在初始化时获取余额
	if Config.Model == openai.GPT3Dot5Turbo0613 || Config.Model == openai.GPT3Dot5Turbo0301 || Config.Model == openai.GPT3Dot5Turbo {
		_, _ = GetBalance()
	}
}
