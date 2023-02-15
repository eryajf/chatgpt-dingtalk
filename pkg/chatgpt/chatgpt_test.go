package chatgpt

import (
	"fmt"
	"testing"
	"time"
)

func TestChatGPT(t *testing.T) {
	chat := New("CHATGPT_API_KEY", "", 0)
	defer chat.Close()

	//select {
	//case <-chat.GetDoneChan():
	//	fmt.Println("time out")
	//}
	question := "你认为2022年世界杯的冠军是谁？\n"
	fmt.Printf("Q: %s\n", question)
	answer, err := chat.Chat(question)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("A: %s\n", answer)

	//Q: 你认为2022年世界杯的冠军是谁？
	//A: 这个问题很难回答，因为2022年世界杯还没有开始，所以没有人知道冠军是谁。

}

func TestChatGPT_ChatWithContext(t *testing.T) {
	chat := New("CHATGPT_API_KEY", "", 10*time.Minute)
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
