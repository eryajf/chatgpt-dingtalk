package public

import (
	"github.com/eryajf/chatgpt-dingtalk/config"
	"github.com/eryajf/chatgpt-dingtalk/pkg/cache"
	"github.com/eryajf/chatgpt-dingtalk/pkg/db"
	"github.com/sashabaranov/go-openai"
)

var UserService cache.UserServiceInterface
var Config *config.Configuration
var Prompt *[]config.Prompt

func InitSvc() {
	// 加载配置
	Config = config.LoadConfig()
	// 加载prompt
	Prompt = config.LoadPrompt()
	// 初始化缓存
	UserService = cache.NewUserService()
	// 初始化数据库
	db.InitDB()
	// 暂时不在初始化时获取余额
	if Config.Model == openai.GPT3Dot5Turbo0301 || Config.Model == openai.GPT3Dot5Turbo {
		_, _ = GetBalance()
	}
}
