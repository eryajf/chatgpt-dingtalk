package public

import (
	"fmt"
	"strings"

	"github.com/eryajf/chatgpt-dingtalk/config"
	"github.com/eryajf/chatgpt-dingtalk/pkg/cache"
	"github.com/eryajf/chatgpt-dingtalk/pkg/dingbot"
	"github.com/eryajf/chatgpt-dingtalk/pkg/logger"
	"github.com/sashabaranov/go-openai"
)

var UserService cache.UserServiceInterface
var Config *config.Configuration
var Prompt *[]config.Prompt

func InitSvc() {
	Config = config.LoadConfig()
	Prompt = config.LoadPrompt()
	UserService = cache.NewUserService()
	if Config.Model == openai.GPT3Dot5Turbo0301 || Config.Model == openai.GPT3Dot5Turbo {
		_, _ = GetBalance()
	}
}

func FirstCheck(rmsg *dingbot.ReceiveMsg) bool {
	lc := UserService.GetUserMode(rmsg.SenderStaffId)
	if lc == "" {
		if Config.DefaultMode == "串聊" {
			return true
		} else {
			return false
		}
	}
	if lc != "" && strings.Contains(lc, "串聊") {
		return true
	}
	return false
}

// ProcessRequest 分析处理请求逻辑
// 主要提供单日请求限额的功能
func CheckRequest(rmsg *dingbot.ReceiveMsg) bool {
	if Config.MaxRequest == 0 {
		return true
	}
	count := UserService.GetUseRequestCount(rmsg.SenderStaffId)
	// 判断访问次数是否超过限制
	if count >= Config.MaxRequest {
		logger.Info(fmt.Sprintf("亲爱的: %s，您今日请求次数已达上限，请明天再来，交互发问资源有限，请务必斟酌您的问题，给您带来不便，敬请谅解!", rmsg.SenderNick))
		_, err := rmsg.ReplyToDingtalk(string(dingbot.TEXT), fmt.Sprintf("一个好的问题，胜过十个好的答案！\n亲爱的: %s，您今日请求次数已达上限，请明天再来，交互发问资源有限，请务必斟酌您的问题，给您带来不便，敬请谅解!", rmsg.SenderNick))
		if err != nil {
			logger.Warning(fmt.Errorf("send message error: %v", err))
		}
		return false
	}
	// 访问次数未超过限制，将计数加1
	UserService.SetUseRequestCount(rmsg.SenderStaffId, count+1)
	return true
}
