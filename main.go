package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

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
		// 去除问题的前后空格
		msgObj.Text.Content = strings.TrimSpace(msgObj.Text.Content)
		// 打印钉钉回调过来的请求明细
		logger.Info(fmt.Sprintf("dingtalk callback parameters: %#v", msgObj))
		// TODO: 校验请求
		if len(msgObj.Text.Content) == 1 || msgObj.Text.Content == "帮助" {
			// 欢迎信息
			_, err := msgObj.ReplyToDingtalk(string(dingbot.MARKDOWN), public.Welcome)
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
				return ship.ErrBadRequest.New(fmt.Errorf("send message error: %v", err))
			}
		} else {
			// 除去帮助之外的逻辑分流在这里处理
			switch {
			case strings.HasPrefix(msgObj.Text.Content, "#图片"):
				return process.ImageGenerate(&msgObj)
			default:
				msgObj.Text.Content, err = process.GeneratePrompt(msgObj.Text.Content)
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

	port := ":" + public.Config.Port
	srv := &http.Server{
		Addr:    port,
		Handler: app,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		logger.Info("🚀 The HTTP Server is running on", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutting down server...")

	// 5秒后强制退出
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", err)
	}
	logger.Info("Server exiting!")
}
