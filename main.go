package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/eryajf/chatgpt-dingtalk/public"
	"github.com/eryajf/chatgpt-dingtalk/public/logger"
	"github.com/eryajf/chatgpt-dingtalk/service"
)

var UserService service.UserServiceInterface

func init() {
	UserService = service.NewUserService()
}

func main() {
	Start()
}

var Welcome string = `Commands:
=================================
🙋 单聊 👉 单独聊天，缺省
🗣 串聊 👉 带上下文聊天
🔃 重置 👉 重置带上下文聊天
🚀 帮助 👉 显示帮助信息
=================================
🚜 例：@我发送 空 或 帮助 将返回此帮助信息
💪 Power By https://github.com/eryajf/chatgpt-dingtalk
`

// 💵 余额 👉 查看接口可调用额度

func Start() {
	// 定义一个处理器函数
	handler := func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Warning("read request body failed: %v\n", err.Error())
			return
		}
		if len(data) == 0 {
			logger.Warning("回调参数为空，以至于无法正常解析，请检查原因")
			return
		}
		var msgObj = new(public.ReceiveMsg)
		err = json.Unmarshal(data, &msgObj)
		if err != nil {
			logger.Warning("unmarshal request body failed: %v\n", err)
		}
		if msgObj.Text.Content == "" || msgObj.ChatbotUserID == "" {
			logger.Warning("从钉钉回调过来的内容为空，根据过往的经验，或许重新创建一下机器人，能解决这个问题")
			return
		}
		// TODO: 校验请求
		if len(msgObj.Text.Content) == 1 || msgObj.Text.Content == " 帮助" {
			// 欢迎信息
			msgObj.ReplyText(Welcome)
		} else {
			logger.Info(fmt.Sprintf("dingtalk callback parameters: %#v", msgObj))
			err = ProcessRequest(*msgObj)
			if err != nil {
				logger.Warning("process request failed: %v\n", err)
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

func FirstCheck(rmsg public.ReceiveMsg) bool {
	lc := UserService.GetUserMode(rmsg.SenderNick)
	if lc != "" && strings.Contains(lc, "串聊") {
		return true
	}
	return false
}

func ProcessRequest(rmsg public.ReceiveMsg) error {
	switch rmsg.Text.Content {
	case " 单聊":
		UserService.SetUserMode(rmsg.SenderNick, rmsg.Text.Content)
		rmsg.ReplyText(fmt.Sprintf("=====现在进入与👉%s👈单聊的模式 =====", rmsg.SenderNick))
	case " 串聊":
		UserService.SetUserMode(rmsg.SenderNick, rmsg.Text.Content)
		rmsg.ReplyText(fmt.Sprintf("=====现在进入与👉%s👈串聊的模式 =====", rmsg.SenderNick))
	case " 重置":
		UserService.ClearUserMode(rmsg.SenderNick)
		err := os.Remove("openaiCache/" + rmsg.SenderNick)
		if err != nil && !strings.Contains(fmt.Sprintf("%s", err), "no such file or directory") {
			rmsg.ReplyText(fmt.Sprintf("=====清理与👉%s👈的对话缓存失败，错误信息: %v\n请检查=====", rmsg.SenderNick, err))
		} else {
			rmsg.ReplyText(fmt.Sprintf("=====已重置与👉%s👈的对话模式，可以开始新的对话=====", rmsg.SenderNick))
		}
	default:
		if FirstCheck(rmsg) {
			cli, reply, err := public.ContextQa(rmsg.Text.Content, rmsg.SenderNick)
			if err != nil {
				logger.Info("gpt request error: %v \n", err)
				_, err = rmsg.ReplyText(fmt.Sprintf("请求openai失败了，错误信息：%v", err))
				if err != nil {
					logger.Warning("send message error: %v \n", err)
					return err
				}
			}
			if reply == "" {
				logger.Warning("get gpt result falied: %v\n", err)
				return nil
			} else {
				reply = strings.TrimSpace(reply)
				reply = strings.Trim(reply, "\n")
				// 回复@我的用户
				replyText := "@" + rmsg.SenderNick + "\n" + reply
				_, err = rmsg.ReplyText(replyText)
				if err != nil {
					logger.Warning("send message error: %v \n", err)
					return err
				}
				path := "openaiCache/" + rmsg.SenderNick
				cli.ChatContext.SaveConversation(path)
			}
		} else {
			reply, err := public.SingleQa(rmsg.Text.Content, rmsg.SenderNick)
			if err != nil {
				logger.Info("gpt request error: %v \n", err)
				_, err = rmsg.ReplyText(fmt.Sprintf("请求openai失败了，错误信息：%v", err))
				if err != nil {
					logger.Warning("send message error: %v \n", err)
					return err
				}
			}
			if reply == "" {
				logger.Warning("get gpt result falied: %v\n", err)
				return nil
			} else {
				reply = strings.TrimSpace(reply)
				reply = strings.Trim(reply, "\n")
				// 回复@我的用户
				replyText := "@" + rmsg.SenderNick + "\n" + reply
				_, err = rmsg.ReplyText(replyText)
				if err != nil {
					logger.Warning("send message error: %v \n", err)
					return err
				}
			}
		}
	}
	return nil
}
