package process

import (
	"fmt"
	"strings"
	"time"

	"github.com/eryajf/chatgpt-dingtalk/pkg/db"
	"github.com/eryajf/chatgpt-dingtalk/pkg/dingbot"
	"github.com/eryajf/chatgpt-dingtalk/pkg/logger"
	"github.com/eryajf/chatgpt-dingtalk/public"
	"github.com/solywsh/chatgpt"
)

// ProcessRequest 分析处理请求逻辑
func ProcessRequest(rmsg *dingbot.ReceiveMsg) error {
	if public.CheckRequest(rmsg) {
		content := strings.TrimSpace(rmsg.Text.Content)
		switch content {
		case "单聊":
			public.UserService.SetUserMode(rmsg.GetSenderIdentifier(), content)
			_, err := rmsg.ReplyToDingtalk(string(dingbot.TEXT), fmt.Sprintf("=====现在进入与👉%s👈单聊的模式 =====", rmsg.SenderNick))
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		case "串聊":
			public.UserService.SetUserMode(rmsg.GetSenderIdentifier(), content)
			_, err := rmsg.ReplyToDingtalk(string(dingbot.TEXT), fmt.Sprintf("=====现在进入与👉%s👈串聊的模式 =====", rmsg.SenderNick))
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		case "重置":
			// 重置用户对话模式
			public.UserService.ClearUserMode(rmsg.GetSenderIdentifier())
			// 清空用户对话上下文
			public.UserService.ClearUserSessionContext(rmsg.GetSenderIdentifier())
			// 清空用户对话的答案ID
			public.UserService.ClearAnswerID(rmsg.SenderNick, rmsg.GetChatTitle())
			_, err := rmsg.ReplyToDingtalk(string(dingbot.TEXT), fmt.Sprintf("=====已重置与👉%s👈的对话模式，可以开始新的对话=====", rmsg.SenderNick))
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		case "模板":
			var title string
			for _, v := range *public.Prompt {
				title = title + v.Title + " | "
			}
			_, err := rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), fmt.Sprintf("%s 您好，当前程序内置集成了这些提示词：\n\n-----\n\n| %s \n\n-----\n\n您可以选择某个提示词作为对话内容的开头。\n\n以周报为例，可发送\"#周报 我本周用Go写了一个钉钉集成ChatGPT的聊天应用\"，可将工作内容填充为一篇完整的周报。\n\n-----\n\n若您不清楚某个提示词的所代表的含义，您可以直接发送提示词，例如直接发送\"#周报\"", rmsg.SenderNick, title))
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		case "图片":
			_, err := rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), "发送以 **#图片** 开头的内容，将会触发绘画能力，图片生成之后，将会保存在程序根目录下的 **images目录** \n 如果你绘图没有思路，可以在这两个网站寻找灵感。\n - [https://lexica.art/](https://lexica.art/)\n- [https://www.clickprompt.org/zh-CN/](https://www.clickprompt.org/zh-CN/)")
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		case "余额":
			// 	cacheMsg := public.UserService.GetUserMode("system_balance")
			// 	if cacheMsg == "" {
			// 		rst, err := public.GetBalance()
			// 		if err != nil {
			// 			logger.Warning(fmt.Errorf("get balance error: %v", err))
			// 			return err
			// 		}
			// 		t1 := time.Unix(int64(rst.Grants.Data[0].EffectiveAt), 0)
			// 		t2 := time.Unix(int64(rst.Grants.Data[0].ExpiresAt), 0)
			// 		cacheMsg = fmt.Sprintf("💵 已用: 💲%v\n💵 剩余: 💲%v\n⏳ 有效时间: 从 %v 到 %v\n", fmt.Sprintf("%.2f", rst.TotalUsed), fmt.Sprintf("%.2f", rst.TotalAvailable), t1.Format("2006-01-02 15:04:05"), t2.Format("2006-01-02 15:04:05"))
			// 	}
			cacheMsg := "官方暂时改写了余额接口，因此暂不提供查询余额功能！2023-04-03"
			_, err := rmsg.ReplyToDingtalk(string(dingbot.TEXT), cacheMsg)
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		case "查对话":
			msg := "使用如下指令进行查询:\n\n---\n\n**#查对话 username:张三**\n\n---\n\n需要注意格式必须严格与上边一致，否则会查询失败\n\n只有钉钉管理员，程序系统管理员，与查自己的情况下，才会被允许"
			_, err := rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), msg)
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
	}
	return nil
}

// 执行处理请求
func Do(mode string, rmsg *dingbot.ReceiveMsg) error {
	// 先把模式注入
	public.UserService.SetUserMode(rmsg.GetSenderIdentifier(), mode)
	switch mode {
	case "单聊":
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
		reply, err := chatgpt.SingleQa(rmsg.Text.Content, rmsg.GetSenderIdentifier())
		if err != nil {
			logger.Info(fmt.Errorf("gpt request error: %v", err))
			if strings.Contains(fmt.Sprintf("%v", err), "maximum text length exceeded") {
				public.UserService.ClearUserSessionContext(rmsg.GetSenderIdentifier())
				_, err = rmsg.ReplyToDingtalk(string(dingbot.TEXT), fmt.Sprintf("请求openai失败了，错误信息：%v，看起来是超过最大对话限制了，已自动重置您的对话", err))
				if err != nil {
					logger.Warning(fmt.Errorf("send message error: %v", err))
					return err
				}
			} else {
				_, err = rmsg.ReplyToDingtalk(string(dingbot.TEXT), fmt.Sprintf("请求openai失败了，错误信息：%v", err))
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
			_, err = rmsg.ReplyToDingtalk(string(dingbot.TEXT), reply)
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
				return err
			}
		}
	case "串聊":
		lastAid := public.UserService.GetAnswerID(rmsg.SenderNick, rmsg.GetChatTitle())
		qObj := db.Chat{
			Username:      rmsg.SenderNick,
			Source:        rmsg.GetChatTitle(),
			ChatType:      db.Q,
			ParentContent: lastAid,
			Content:       rmsg.Text.Content,
		}
		qid, err := qObj.Add()
		if err != nil {
			logger.Error("往MySQL新增数据失败,错误信息：", err)
		}
		cli, reply, err := chatgpt.ContextQa(rmsg.Text.Content, rmsg.GetSenderIdentifier())
		if err != nil {
			logger.Info(fmt.Sprintf("gpt request error: %v", err))
			if strings.Contains(fmt.Sprintf("%v", err), "maximum text length exceeded") {
				public.UserService.ClearUserSessionContext(rmsg.GetSenderIdentifier())
				_, err = rmsg.ReplyToDingtalk(string(dingbot.TEXT), fmt.Sprintf("请求openai失败了，错误信息：%v，看起来是超过最大对话限制了，已自动重置您的对话", err))
				if err != nil {
					logger.Warning(fmt.Errorf("send message error: %v", err))
					return err
				}
			} else {
				_, err = rmsg.ReplyToDingtalk(string(dingbot.TEXT), fmt.Sprintf("请求openai失败了，错误信息：%v", err))
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
			aObj := db.Chat{
				Username:      rmsg.SenderNick,
				Source:        rmsg.GetChatTitle(),
				ChatType:      db.A,
				ParentContent: qid,
				Content:       reply,
			}
			aid, err := aObj.Add()
			if err != nil {
				logger.Error("往MySQL新增数据失败,错误信息：", err)
			}
			// 将当前回答的ID放入缓存
			public.UserService.SetAnswerID(rmsg.SenderNick, rmsg.GetChatTitle(), aid)
			logger.Info(fmt.Sprintf("🤖 %s得到的答案: %#v", rmsg.SenderNick, reply))
			// 回复@我的用户
			_, err = rmsg.ReplyToDingtalk(string(dingbot.TEXT), reply)
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
				return err
			}
			_ = cli.ChatContext.SaveConversation(rmsg.GetSenderIdentifier())
		}
	default:

	}
	return nil
}

func ImageGenerate(rmsg *dingbot.ReceiveMsg) error {
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
	reply, err := chatgpt.ImageQa(rmsg.Text.Content, rmsg.GetSenderIdentifier())
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
func SelectHistory(rmsg *dingbot.ReceiveMsg) error {
	name := strings.TrimSpace(strings.Split(rmsg.Text.Content, ":")[1])
	if !rmsg.IsAdmin && name != rmsg.SenderNick && !public.JudgeAdminUsers(rmsg.SenderNick) {
		_, err := rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), "**🤷 抱歉，您没有权限查询其他人的对话记录！**")
		if err != nil {
			logger.Error(fmt.Errorf("send message error: %v", err))
			return err
		}
		return nil
	}
	// 获取数据列表
	var chat db.Chat
	if !chat.Exist(map[string]interface{}{"username": name}) {
		_, err := rmsg.ReplyToDingtalk(string(dingbot.TEXT), "用户名错误，这个用户不存在，请核实之后再进行查询")
		if err != nil {
			logger.Error(fmt.Errorf("send message error: %v", err))
			return err
		}
		return fmt.Errorf("用户名错误，这个用户不存在，请核实之后重新查询")
	}
	chats, err := chat.List(db.ChatListReq{
		Username: name,
	})
	if err != nil {
		return err
	}
	var rst string
	for _, chatTmp := range chats {
		ctime := chatTmp.CreatedAt.Format("2006-01-02 15:04:05")
		if chatTmp.ChatType == 1 {
			rst += fmt.Sprintf("## 🙋 %s 问\n\n**时间:** %v\n\n**问题为:** %s\n\n", chatTmp.Username, ctime, chatTmp.Content)
		} else {
			rst += fmt.Sprintf("## 🤖 机器人答\n\n**时间:** %v\n\n**回答如下：** \n\n%s\n\n", ctime, chatTmp.Content)
		}
		// TODO: 答案应该严格放在问题之后，目前只根据ID排序进行的陈列，当一个用户同时提出多个问题时，最终展示的可能会有点问题
	}
	fileName := time.Now().Format("20060102-150405") + ".md"
	// 写入文件
	if err = public.WriteToFile("./data/chatHistory/"+fileName, []byte(rst)); err != nil {
		return err
	}
	// 回复@我的用户
	reply := fmt.Sprintf("- 在线查看: [点我](%s)\n- 下载文件: [点我](%s)\n- 在线预览请安装插件:[Markdown Preview Plus](https://chrome.google.com/webstore/detail/markdown-preview-plus/febilkbfcbhebfnokafefeacimjdckgl)", public.Config.ServiceURL+"/history/"+fileName, public.Config.ServiceURL+"/download/"+fileName)
	logger.Info(fmt.Sprintf("🤖 %s得到的答案: %#v", rmsg.SenderNick, reply))
	_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), reply)
	if err != nil {
		logger.Error(fmt.Errorf("send message error: %v", err))
		return err
	}
	return nil
}
