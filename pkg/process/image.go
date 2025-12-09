package process

import (
	"context"
	"fmt"
	"strings"

	"github.com/eryajf/chatgpt-dingtalk/pkg/db"
	"github.com/eryajf/chatgpt-dingtalk/pkg/dingbot"
	"github.com/eryajf/chatgpt-dingtalk/pkg/llm"
	"github.com/eryajf/chatgpt-dingtalk/pkg/logger"
	"github.com/eryajf/chatgpt-dingtalk/public"
)

// ImageGenerate openai生成图片
func ImageGenerate(ctx context.Context, rmsg *dingbot.ReceiveMsg) error {
	if public.Config.AzureOn {
		_, err := rmsg.ReplyToDingtalk(string(dingbot.
			MARKDOWN), "azure 模式下暂不支持图片创作功能")
		if err != nil {
			logger.Warning(fmt.Errorf("send message error: %v", err))
		}
		return err
	}
	qObj := db.Chat{
		Username:      rmsg.SenderNick,
		Source:        rmsg.GetChatTitle(),
		ChatType:      db.Q,
		ParentContent: 0,
		Content:       rmsg.Text.Content,
	}
	qid, err := qObj.Add()
	if err != nil {
		logger.Error("往MySQL新增数据失败,错误信息：", err)
	}
	reply, err := llm.ImageQa(ctx, rmsg.Text.Content, rmsg.GetSenderIdentifier())
	if err != nil {
		logger.Info(fmt.Errorf("gpt request error: %v", err))
		_, err = rmsg.ReplyToDingtalk(string(dingbot.TEXT), fmt.Sprintf("请求openai失败了，错误信息：%v", err))
		if err != nil {
			logger.Error(fmt.Errorf("send message error: %v", err))
			return err
		}
	}
	if reply == "" {
		logger.Warning(fmt.Errorf("get gpt result falied: %v", err))
		return nil
	} else {
		reply = strings.TrimSpace(reply)
		reply = strings.Trim(reply, "\n")
		reply = fmt.Sprintf(">点击图片可旋转或放大。\n![](%s)", reply)
		aObj := db.Chat{
			Username:      rmsg.SenderNick,
			Source:        rmsg.GetChatTitle(),
			ChatType:      db.A,
			ParentContent: qid,
			Content:       reply,
		}
		_, err := aObj.Add()
		if err != nil {
			logger.Error("往MySQL新增数据失败,错误信息：", err)
		}
		logger.Info(fmt.Sprintf("🤖 %s得到的答案: %#v", rmsg.SenderNick, reply))
		// 回复@我的用户
		_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), reply)
		if err != nil {
			logger.Error(fmt.Errorf("send message error: %v", err))
			return err
		}
	}
	return nil
}
