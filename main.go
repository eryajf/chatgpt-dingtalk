package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/eryajf/chatgpt-dingtalk/config"
	"github.com/eryajf/chatgpt-dingtalk/public"
	"github.com/eryajf/chatgpt-dingtalk/public/logger"
	"github.com/solywsh/chatgpt"
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
			logger.Info(fmt.Sprintf("dingtalk callback parameters: %#v", msgObj))
			err = ProcessRequest(*msgObj)
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

func ProcessRequest(rmsg public.ReceiveMsg) error {
	content := strings.TrimSpace(rmsg.Text.Content)
	switch content {
	case "单聊":
		public.UserService.SetUserMode(rmsg.SenderStaffId, content)
		_, err := rmsg.ReplyText(fmt.Sprintf("=====现在进入与👉%s👈单聊的模式 =====", rmsg.SenderNick), rmsg.SenderStaffId)
		if err != nil {
			logger.Warning(fmt.Errorf("send message error: %v", err))
		}
	case "串聊":
		public.UserService.SetUserMode(rmsg.SenderStaffId, content)
		_, err := rmsg.ReplyText(fmt.Sprintf("=====现在进入与👉%s👈串聊的模式 =====", rmsg.SenderNick), rmsg.SenderStaffId)
		if err != nil {
			logger.Warning(fmt.Errorf("send message error: %v", err))
		}
	case "重置":
		public.UserService.ClearUserMode(rmsg.SenderStaffId)
		public.UserService.ClearUserSessionContext(rmsg.SenderStaffId)
		_, err := rmsg.ReplyText(fmt.Sprintf("=====已重置与👉%s👈的对话模式，可以开始新的对话=====", rmsg.SenderNick), rmsg.SenderStaffId)
		if err != nil {
			logger.Warning(fmt.Errorf("send message error: %v", err))
		}
	case "余额":
		cacheMsg := public.UserService.GetUserMode("system_balance")
		if cacheMsg == "" {
			rst, err := public.GetBalance()
			if err != nil {
				logger.Warning(fmt.Errorf("get balance error: %v", err))
				return err
			}
			t1 := time.Unix(int64(rst.Grants.Data[0].EffectiveAt), 0)
			t2 := time.Unix(int64(rst.Grants.Data[0].ExpiresAt), 0)
			cacheMsg = fmt.Sprintf("💵 已用: 💲%v\n💵 剩余: 💲%v\n⏳ 有效时间: 从 %v 到 %v\n", fmt.Sprintf("%.2f", rst.TotalUsed), fmt.Sprintf("%.2f", rst.TotalAvailable), t1.Format("2006-01-02 15:04:05"), t2.Format("2006-01-02 15:04:05"))
		}

		_, err := rmsg.ReplyText(cacheMsg, rmsg.SenderStaffId)
		if err != nil {
			logger.Warning(fmt.Errorf("send message error: %v", err))
		}
	default:
		if public.FirstCheck(rmsg) {
			return Do("串聊", rmsg)
		} else {
			return Do("单聊", rmsg)
		}
	}
	return nil
}

func Do(mode string, rmsg public.ReceiveMsg) error {
	// 先把模式注入
	public.UserService.SetUserMode(rmsg.SenderStaffId, mode)
	switch mode {
	case "单聊":
		reply, err := SingleQa(rmsg.Text.Content, rmsg.SenderStaffId)
		if err != nil {
			logger.Info(fmt.Errorf("gpt request error: %v", err))
			if strings.Contains(fmt.Sprintf("%v", err), "maximum text length exceeded") {
				public.UserService.ClearUserSessionContext(rmsg.SenderStaffId)
				_, err = rmsg.ReplyText(fmt.Sprintf("请求openai失败了，错误信息：%v，看起来是超过最大对话限制了，已自动重置您的对话", err), rmsg.SenderStaffId)
				if err != nil {
					logger.Warning(fmt.Errorf("send message error: %v", err))
					return err
				}
			} else {
				_, err = rmsg.ReplyText(fmt.Sprintf("请求openai失败了，错误信息：%v", err), rmsg.SenderStaffId)
				if err != nil {
					logger.Warning(fmt.Errorf("send message error: %v", err))
					return err
				}
			}
		}
		if reply == "" {
			logger.Warning(fmt.Errorf("get gpt result falied: %v", err))
			return nil
		} else {
			reply = strings.TrimSpace(reply)
			reply = strings.Trim(reply, "\n")
			// 回复@我的用户
			// fmt.Println("单聊结果是：", reply)
			_, err = rmsg.ReplyText(reply, rmsg.SenderStaffId)
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
				return err
			}
		}
	case "串聊":
		cli, reply, err := ContextQa(rmsg.Text.Content, rmsg.SenderStaffId)
		if err != nil {
			logger.Info(fmt.Sprintf("gpt request error: %v", err))
			if strings.Contains(fmt.Sprintf("%v", err), "maximum text length exceeded") {
				public.UserService.ClearUserSessionContext(rmsg.SenderStaffId)
				_, err = rmsg.ReplyText(fmt.Sprintf("请求openai失败了，错误信息：%v，看起来是超过最大对话限制了，已自动重置您的对话", err), rmsg.SenderStaffId)
				if err != nil {
					logger.Warning(fmt.Errorf("send message error: %v", err))
					return err
				}
			} else {
				_, err = rmsg.ReplyText(fmt.Sprintf("请求openai失败了，错误信息：%v", err), rmsg.SenderStaffId)
				if err != nil {
					logger.Warning(fmt.Errorf("send message error: %v", err))
					return err
				}
			}
		}
		if reply == "" {
			logger.Warning(fmt.Errorf("get gpt result falied: %v", err))
			return nil
		} else {
			reply = strings.TrimSpace(reply)
			reply = strings.Trim(reply, "\n")
			// 回复@我的用户
			_, err = rmsg.ReplyText(reply, rmsg.SenderStaffId)
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
				return err
			}
			_ = cli.ChatContext.SaveConversation(rmsg.SenderStaffId)
		}
	default:

	}
	return nil
}

func SingleQa(question, userId string) (answer string, err error) {
	cfg := config.LoadConfig()
	chat := chatgpt.New(cfg.ApiKey, cfg.HttpProxy, userId, cfg.SessionTimeout)
	defer chat.Close()
	return chat.ChatWithContext(question)
}

func ContextQa(question, userId string) (chat *chatgpt.ChatGPT, answer string, err error) {
	cfg := config.LoadConfig()
	chat = chatgpt.New(cfg.ApiKey, cfg.HttpProxy, userId, cfg.SessionTimeout)
	if public.UserService.GetUserSessionContext(userId) != "" {
		err = chat.ChatContext.LoadConversation(userId)
		if err != nil {
			fmt.Printf("load station failed: %v\n", err)
		}
	}
	answer, err = chat.ChatWithContext(question)
	return
}
