package process

import (
	"strings"

	"github.com/eryajf/chatgpt-dingtalk/public"
)

// GeneratePrompt 生成当次请求的 Prompt
func GeneratePrompt(msg string) (rst string) {
	for _, prompt := range *public.Prompt {
		if strings.HasPrefix(msg, prompt.Title) {
			rst = prompt.Content + strings.Replace(msg, prompt.Title, "", -1)
			return
		} else {
			rst = msg
		}
	}
	return
}
