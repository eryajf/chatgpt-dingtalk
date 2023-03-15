package chatgpt

import (
	"os"
	"testing"
)

func TestOfflineContext(t *testing.T) {
	key := os.Getenv("CHATGPT_API_KEY")
	if key == "" {
		t.Skip("CHATGPT_API_KEY is not set")
	}
	cli := New("")
	reply, err := cli.ChatWithContext("我叫老三，你是？")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("我叫老三，你是？ => %s", reply)

	err = cli.ChatContext.SaveConversation("test.conversation")
	if err != nil {
		t.Fatalf("储存对话记录失败: %v", err)
	}
	cli.ChatContext.ResetConversation("")

	reply, err = cli.ChatWithContext("你知道我是谁吗?")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("你知道我是谁吗? => %s", reply)
	// assert.NotContains(t, reply, "老三")

	err = cli.ChatContext.LoadConversation("test.conversation")
	if err != nil {
		t.Fatalf("读取对话记录失败: %v", err)
	}

	reply, err = cli.ChatWithContext("你知道我是谁吗?")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("你知道我是谁吗? => %s", reply)

	// AI 理应知道他叫老三
	// assert.Contains(t, reply, "老三")
}

func TestMaintainContext(t *testing.T) {
	key := os.Getenv("CHATGPT_API_KEY")
	if key == "" {
		t.Skip("CHATGPT_API_KEY is not set")
	}
	cli := New("")
	cli.ChatContext = NewContext(
		WithMaxSeqTimes(1),
		WithMaintainSeqTimes(true),
	)

	reply, err := cli.ChatWithContext("我叫老三，你是？")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("我叫老三，你是？ => %s", reply)

	reply, err = cli.ChatWithContext("你知道我是谁吗?")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("你知道我是谁吗? => %s", reply)

	// 对话次数已经超过 1 次，因此最先前的对话已被移除，AI 理应不知道他叫老三
	// assert.NotContains(t, reply, "老三")
}

func init() {
	// 本地加载适用于本地测试，如果要在github进行测试，可以透过传入 secrets 到环境参数
	// _ = godotenv.Load(".env.local")
}
