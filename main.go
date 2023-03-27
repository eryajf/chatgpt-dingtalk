package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/eryajf/chatgpt-dingtalk/pkg/dingbot"
	"github.com/eryajf/chatgpt-dingtalk/pkg/logger"
	"github.com/eryajf/chatgpt-dingtalk/pkg/process"
	"github.com/eryajf/chatgpt-dingtalk/public"
	"github.com/xgfone/ship/v5"
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
🎨 图片 👉 根据prompt生成图片
=================================
🚜 例：@我发送 空 或 帮助 将返回此帮助信息
💪 Power By https://github.com/eryajf/chatgpt-dingtalk
`

func Start() {
	app := ship.Default()
	app.Route("/").POST(func(c *ship.Context) error {
		var msgObj dingbot.ReceiveMsg
		err := c.Bind(&msgObj)
		if err != nil {
			return ship.ErrBadRequest.New(fmt.Errorf("bind to receivemsg failed : %v", err))
		}
		if msgObj.Text.Content == "" || msgObj.ChatbotUserID == "" {
			logger.Warning("从钉钉回调过来的内容为空，根据过往的经验，或许重新创建一下机器人，能解决这个问题")
			return ship.ErrBadRequest.New(fmt.Errorf("从钉钉回调过来的内容为空，根据过往的经验，或许重新创建一下机器人，能解决这个问题"))
		}

		// 打印钉钉回调过来的请求明细
		logger.Info(fmt.Sprintf("dingtalk callback parameters: %#v", msgObj))
		// TODO: 校验请求
		if len(msgObj.Text.Content) == 1 || strings.TrimSpace(msgObj.Text.Content) == "帮助" {
			// 欢迎信息
			_, err := msgObj.ReplyToDingtalk(string(dingbot.TEXT), Welcome)
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
				return ship.ErrBadRequest.New(fmt.Errorf("send message error: %v", err))
			}
		} else {
			// 除去帮助之外的逻辑分流在这里处理
			switch {
			case strings.HasPrefix(strings.TrimSpace(msgObj.Text.Content), "#图片"):
				return process.ImageGenerate(&msgObj)
			default:
				msgObj.Text.Content = process.GeneratePrompt(strings.TrimSpace(msgObj.Text.Content))
				logger.Info(fmt.Sprintf("after generate prompt: %#v", msgObj.Text.Content))
				return process.ProcessRequest(&msgObj)
			}
		}
		return nil
	})
	// 解析生成后的图片
	app.Route("/images/:filename").GET(func(c *ship.Context) error {
		filename := c.Param("filename")
		root := "./images/"
		return c.File(filepath.Join(root, filename))
	})

	// 服务端口
	port := ":" + public.Config.Port
	// 启动服务器
	ship.StartServer(port, app)
}
