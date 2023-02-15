package public

import (
	"fmt"

	"github.com/eryajf/chatgpt-dingtalk/config"
	"github.com/solywsh/chatgpt"
)

func SingleQa(question, userId string) (answer string, err error) {
	cfg := config.LoadConfig()
	chat := chatgpt.New(cfg.ApiKey, userId, cfg.SessionTimeout)
	defer chat.Close()
	return chat.ChatWithContext(question)
}

func ContextQa(question, userId string) (chat *chatgpt.ChatGPT, answer string, err error) {
	cfg := config.LoadConfig()
	chat = chatgpt.New(cfg.ApiKey, userId, cfg.SessionTimeout)
	path := "openaiCache/" + userId
	err = chat.ChatContext.LoadConversation(path)
	if err != nil {
		fmt.Printf("load station failed: %v\n", err)
	}
	answer, err = chat.ChatWithContext(question)
	return
}
