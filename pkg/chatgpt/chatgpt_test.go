package chatgpt

import (
	"fmt"
	"github.com/eryajf/chatgpt-dingtalk/public"
	"testing"
)

func TestChatGPT_ChatWithContext(t *testing.T) {
	public.InitSvc()
	chat := New("")
	defer chat.Close()
	//go func() {
	//	select {
	//	case <-chat.GetDoneChan():
	//		fmt.Println("time out")
	//	}
	//}()
	question := "现在你是一只猫，接下来你只能用\"喵喵喵\"回答."
	fmt.Printf("Q: %s\n", question)
	answer, err := chat.ChatWithContext(question)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("A: %s\n", answer)
	question = "你是一只猫吗？"
	fmt.Printf("Q: %s\n", question)
	answer, err = chat.ChatWithContext(question)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("A: %s\n", answer)
}
