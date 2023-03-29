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

var Welcome string = `# 发送信息

若您想给机器人发送信息，请选择：

1. 在本机器人所在群里@机器人；
2. 点击机器人的头像后，再点击"发消息"。

机器人收到您的信息后，默认会交给chatgpt进行处理。除非，您发送的内容是7个**系统指令**之一。

-----

# 系统指令

系统指令是一些特殊的词语，当您向机器人发送这些词语时，会触发对应的功能：

**单聊**：每条消息都是单独的对话，不包含上下文

**串聊**：对话会携带上下文，除非您主动重置对话或对话长度超过限制

**重置**：重置上下文

**余额**：查询机器人所用OpenAI账号的余额

**模板**：查询机器人内置的快捷模板

**图片**：查看如何根据提示词生成图片

**帮助**：重新获取帮助信息

-----

# 友情提示

使用"串聊模式"会显著加快机器人所用账号的余额消耗速度。

因此，若无保留上下文的需求，建议使用"单聊模式"。

即使有保留上下文的需求，也应适时使用"重置"指令来重置上下文。

-----

# 项目地址

本项目已在GitHub开源，[查看源代码](https://github.com/eryajf/chatgpt-dingtalk)。
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
			_, err := msgObj.ReplyToDingtalk(string(dingbot.MARKDOWN), Welcome)
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
				msgObj.Text.Content, err = process.GeneratePrompt(strings.TrimSpace(msgObj.Text.Content))
				// err不为空：提示词之后没有文本 -> 直接返回提示词所代表的内容
				if err != nil {
					_, err = msgObj.ReplyToDingtalk(string(dingbot.TEXT), msgObj.Text.Content)
					if err != nil {
						logger.Warning(fmt.Errorf("send message error: %v", err))
						return err
					}
					return nil
				}
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
