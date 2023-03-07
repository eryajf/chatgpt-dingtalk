package public

import (
	"strings"

	"github.com/eryajf/chatgpt-dingtalk/config"
	"github.com/eryajf/chatgpt-dingtalk/service"
)

var UserService service.UserServiceInterface
var Config *config.Configuration

func InitSvc() {
	Config = config.LoadConfig()
	UserService = service.NewUserService()
	_, _ = GetBalance()
}

func FirstCheck(rmsg ReceiveMsg) bool {
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
