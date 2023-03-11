package public

import (
	"strings"

	"github.com/eryajf/chatgpt-dingtalk/config"
	"github.com/eryajf/chatgpt-dingtalk/service"
)

var UserService service.UserServiceInterface
var Config *config.Configuration
var ApiKeyList *ApiKeyInfoList

func InitSvc() {
	Config = config.LoadConfig()
	UserService = service.NewUserService()
	ApiKeyList = InitApiKeyInfo()
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

func CheckAllowGroups(rmsg ReceiveMsg) bool {
	if len(Config.AllowGroups) == 0 {
		return true
	}

	for _, v := range Config.AllowGroups {
		if rmsg.ConversationTitle == v {
			return true
		}
	}
	return false
}

func CheckAllowUsers(rmsg ReceiveMsg) bool {
	if len(Config.AllowUsers) == 0 {
		return true
	}

	for _, v := range Config.AllowUsers {
		if rmsg.SenderNick == v {
			return true
		}
	}
	return false
}
