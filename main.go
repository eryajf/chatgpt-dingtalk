package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/eryajf/chatgpt-dingtalk/gtp"
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

func Start() {
	// 定义一个处理器函数
	handler := func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Warning("read request body failed: %v\n", err.Error())
			return
		}
		// TODO: 校验请求
		// fmt.Println(r.Header)
		if len(data) == 0 {
			logger.Warning("回调参数为空,以至于无法正常解析,请检查原因")
			return
		} else {
			var msgObj = new(public.ReceiveMsg)
			err = json.Unmarshal(data, &msgObj)
			if err != nil {
				logger.Warning("unmarshal request body failed: %v\n", err)
			}
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

func ProcessRequest(rmsg public.ReceiveMsg) error {
	atText := "@" + rmsg.SenderNick + "\n" + " "
	if UserService.ClearUserSessionContext(rmsg.SenderID, rmsg.Text.Content) {
		_, err := rmsg.ReplyText(atText + "上下文已经清空了，你可以问下一个问题啦。")
		if err != nil {
			logger.Warning("response user error: %v \n", err)
			return err
		}
	} else {
		requestText := getRequestText(rmsg)
		// 获取问题的答案
		reply, err := gtp.Completions(requestText)
		if err != nil {
			logger.Info("gtp request error: %v \n", err)
			_, err = rmsg.ReplyText("机器人太累了，让她休息会儿，过一会儿再来请求。")
			if err != nil {
				logger.Warning("send message error: %v \n", err)
				return err
			}
			logger.Info("request openai error: %v\n", err)
			return err
		}
		if reply == "" {
			logger.Warning("get gpt result falied: %v\n", err)
			return nil
		}
		// 回复@我的用户
		reply = strings.TrimSpace(reply)
		reply = strings.Trim(reply, "\n")

		UserService.SetUserSessionContext(rmsg.SenderID, requestText, reply)
		replyText := atText + reply
		_, err = rmsg.ReplyText(replyText)
		if err != nil {
			logger.Info("send message error: %v \n", err)
			return err
		}
	}
	return nil
}

// getRequestText 获取请求接口的文本，要做一些清洗
func getRequestText(rmsg public.ReceiveMsg) string {
	// 1.去除空格以及换行
	requestText := strings.TrimSpace(rmsg.Text.Content)
	requestText = strings.Trim(rmsg.Text.Content, "\n")

	// 2.替换掉当前用户名称
	replaceText := "@" + rmsg.SenderNick
	requestText = strings.TrimSpace(strings.ReplaceAll(rmsg.Text.Content, replaceText, ""))
	if requestText == "" {
		return ""
	}

	// 3.获取上下文，拼接在一起，如果字符长度超出4000，截取为4000。（GPT按字符长度算）
	requestText = UserService.GetUserSessionContext(rmsg.SenderID) + requestText
	if len(requestText) >= 4000 {
		requestText = requestText[:4000]
	}

	// 4.返回请求文本
	return requestText
}
