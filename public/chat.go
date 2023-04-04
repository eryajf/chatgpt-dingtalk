package public

import (
	"strings"

	"github.com/eryajf/chatgpt-dingtalk/pkg/dingbot"
)

func FirstCheck(rmsg *dingbot.ReceiveMsg) bool {
	lc := UserService.GetUserMode(rmsg.GetSenderIdentifier())
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
