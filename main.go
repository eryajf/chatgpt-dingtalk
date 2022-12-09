package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/eryajf/chatgpt-dingtalk/gtp"
	"github.com/eryajf/chatgpt-dingtalk/public"
	"github.com/eryajf/chatgpt-dingtalk/service"
)

var UserService service.UserServiceInterface

func init() {
	UserService = service.NewUserService()
}

func main() {
	// 定义一个处理器函数
	handler := func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// TODO: 校验请求
		// fmt.Println(r.Header)

		var msgObj = new(public.ReceiveMsg)
		err = json.Unmarshal(data, &msgObj)
		if err != nil {
			log.Printf("unmarshal request body failed: %v\n", err)
		}
		err = ProcessRequest(*msgObj)
		if err != nil {
			log.Printf("process request failed: %v\n", err)
		}

	}

	// 创建一个新的 HTTP 服务器
	server := &http.Server{
		Addr:    ":8090",
		Handler: http.HandlerFunc(handler),
	}

	// 启动服务器
	log.Print("Start Listen On ", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func ProcessRequest(rmsg public.ReceiveMsg) error {
	// 获取问题的答案
	reply, err := gtp.Completions(rmsg.Text.Content)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		_, err = rmsg.ReplyText("机器人太累了，让她休息会儿，过一会儿再来请求。")
		if err != nil {
			log.Printf("send message error: %v \n", err)
			return err
		}
		log.Printf("request openai error: %v\n", err)
		return err
	}
	if reply == "" {
		return nil
	}
	// 回复@我的用户
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	atText := "@" + rmsg.SenderNick + "\n" + " "
	// 设置上下文
	if UserService.ClearUserSessionContext(rmsg.SenderID, rmsg.Text.Content) {
		_, err = rmsg.ReplyText(atText + "上下文已经清空了，你可以问下一个问题啦。")
		if err != nil {
			log.Printf("response user error: %v \n", err)
			return err
		}
	}
	UserService.SetUserSessionContext(rmsg.SenderID, rmsg.Text.Content, reply)
	replyText := atText + reply
	_, err = rmsg.ReplyText(replyText)
	if err != nil {
		log.Printf("send message error: %v \n", err)
		return err
	}
	return nil
}
