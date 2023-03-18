package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/eryajf/chatgpt-dingtalk/pkg/process"
	"github.com/eryajf/chatgpt-dingtalk/public"
	"github.com/eryajf/chatgpt-dingtalk/public/logger"
)

func init() {
	public.InitSvc()
}
func main() {
	Start()
}

var Welcome string = `Commands:
=================================
🙋 单聊 👉 单独聊天
📣 串聊 👉 带上下文聊天
🔃 重置 👉 重置带上下文聊天
💵 余额 👉 查询剩余额度
🚀 帮助 👉 显示帮助信息
🌈 模板 👉 内置的prompt
=================================
🚜 例：@我发送 空 或 帮助 将返回此帮助信息
💪 Power By https://github.com/eryajf/chatgpt-dingtalk
`

func Start() {
	// 定义一个处理器函数
	handler := func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Warning(fmt.Sprintf("read request body failed: %v\n", err.Error()))
			return
		}
		if len(data) == 0 {
			logger.Warning("回调参数为空，以至于无法正常解析，请检查原因")
			return
		}
		var msgObj = new(public.ReceiveMsg)
		err = json.Unmarshal(data, &msgObj)
		if err != nil {
			logger.Warning(fmt.Errorf("unmarshal request body failed: %v", err))
		}
		if msgObj.Text.Content == "" || msgObj.ChatbotUserID == "" {
			logger.Warning("从钉钉回调过来的内容为空，根据过往的经验，或许重新创建一下机器人，能解决这个问题")
			return
		}

		// TODO: 校验请求
		if len(msgObj.Text.Content) == 1 || strings.TrimSpace(msgObj.Text.Content) == "帮助" {
			// 欢迎信息
			_, err := msgObj.ReplyText(Welcome, msgObj.SenderStaffId)
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		} else {
			msgObj.Text.Content = process.GeneratePrompt(msgObj.Text.Content)
			logger.Info(fmt.Sprintf("dingtalk callback parameters: %#v", msgObj))
			err = process.ProcessRequest(*msgObj)
			if err != nil {
				logger.Warning(fmt.Errorf("process request failed: %v", err))
			}
		}
	}

	// 创建一个新的 HTTP 服务器
	server := &http.Server{
		Addr:    ":8090",
		Handler: http.HandlerFunc(handler),
	}

	// 启动服务器
	logger.Info("Start Listen On ", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		logger.Danger(err)
	}
}
