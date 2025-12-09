package process

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/eryajf/chatgpt-dingtalk/pkg/db"
	"github.com/eryajf/chatgpt-dingtalk/pkg/dingbot"
	"github.com/eryajf/chatgpt-dingtalk/pkg/llm"
	"github.com/eryajf/chatgpt-dingtalk/pkg/logger"
	"github.com/eryajf/chatgpt-dingtalk/public"
)

// ProcessRequest 分析处理请求逻辑
func ProcessRequest(rmsg *dingbot.ReceiveMsg) error {
	if CheckRequestTimes(rmsg) {
		content := strings.TrimSpace(rmsg.Text.Content)
		timeoutStr := ""
		if content != public.Config.DefaultMode {
			timeoutStr = fmt.Sprintf("\n\n>%s 后将恢复默认聊天模式：%s", FormatTimeDuation(public.Config.SessionTimeout), public.Config.DefaultMode)
		}
		switch content {
		case "单聊":
			public.UserService.SetUserMode(rmsg.GetSenderIdentifier(), content)
			_, err := rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), fmt.Sprintf("**[Concentrate] 现在进入与 %s 的单聊模式**%s", rmsg.SenderNick, timeoutStr))
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		case "串聊":
			public.UserService.SetUserMode(rmsg.GetSenderIdentifier(), content)
			_, err := rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), fmt.Sprintf("**[Concentrate] 现在进入与 %s 的串聊模式**%s", rmsg.SenderNick, timeoutStr))
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		case "重置", "退出", "结束":
			// 重置用户对话模式
			public.UserService.ClearUserMode(rmsg.GetSenderIdentifier())
			// 清空用户对话上下文
			public.UserService.ClearUserSessionContext(rmsg.GetSenderIdentifier())
			// 清空用户对话的答案ID
			public.UserService.ClearAnswerID(rmsg.SenderNick, rmsg.GetChatTitle())
			_, err := rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), fmt.Sprintf("[RecyclingSymbol]已重置与**%s** 的对话模式\n\n> 可以开始新的对话 [Bubble]", rmsg.SenderNick))
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
			if public.Config.AzureOn {
				_, err := rmsg.ReplyToDingtalk(string(dingbot.
					MARKDOWN), "azure 模式下暂不支持图片创作功能")
				if err != nil {
					logger.Warning(fmt.Errorf("send message error: %v", err))
				}
				return err
			}
			_, err := rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), "发送以 **#图片** 开头的内容，将会触发绘画能力，图片生成之后，将会通过消息回复给您。建议尽可能描述需要生成的图片内容及相关细节。\n 如果你绘图没有思路，可以在这两个网站寻找灵感。\n - [https://lexica.art/](https://lexica.art/)\n- [https://www.clickprompt.org/zh-CN/](https://www.clickprompt.org/zh-CN/)")
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		case "余额":
			if public.JudgeAdminUsers(rmsg.SenderStaffId) {
				cacheMsg := public.UserService.GetUserMode("system_balance")
				if cacheMsg == "" {
					rst, err := public.GetBalance()
					if err != nil {
						logger.Warning(fmt.Errorf("get balance error: %v", err))
						return err
					}
					cacheMsg = rst
				}
				_, err := rmsg.ReplyToDingtalk(string(dingbot.TEXT), cacheMsg)
				if err != nil {
					logger.Warning(fmt.Errorf("send message error: %v", err))
				}
			}
		case "查对话":
			if public.JudgeAdminUsers(rmsg.SenderStaffId) {
				msg := "使用如下指令进行查询:\n\n---\n\n**#查对话 username:张三**\n\n---\n\n需要注意格式必须严格与上边一致，否则将会查询失败\n\n只有程序系统管理员有权限查询，即config.yml中的admin_users指定的人员。"
				_, err := rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), msg)
				if err != nil {
					logger.Warning(fmt.Errorf("send message error: %v", err))
				}
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
		reply, err := llm.SingleQa(rmsg.Text.Content, rmsg.GetSenderIdentifier())
		if err != nil {
			logger.Info(fmt.Errorf("gpt request error: %v", err))
			if strings.Contains(fmt.Sprintf("%v", err), "maximum question length exceeded") {
				public.UserService.ClearUserSessionContext(rmsg.GetSenderIdentifier())
				_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), fmt.Sprintf("[Wrong] 请求 OpenAI 失败了\n\n> 错误信息:%v\n\n> 已超过最大文本限制，请缩短提问文字的字数。", err))
				if err != nil {
					logger.Warning(fmt.Errorf("send message error: %v", err))
					return err
				}
			} else {
				_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), fmt.Sprintf("[Wrong] 请求 OpenAI 失败了\n\n> 错误信息:%v", err))
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
			if public.JudgeSensitiveWord(reply) {
				reply = public.SolveSensitiveWord(reply)
			}
			// 回复@我的用户
			_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), FormatMarkdown(reply))
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
		cli, reply, err := llm.ContextQa(rmsg.Text.Content, rmsg.GetSenderIdentifier())
		if err != nil {
			logger.Info(fmt.Sprintf("gpt request error: %v", err))
			if strings.Contains(fmt.Sprintf("%v", err), "maximum text length exceeded") {
				public.UserService.ClearUserSessionContext(rmsg.GetSenderIdentifier())
				_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), fmt.Sprintf("[Wrong] 请求 OpenAI 失败了\n\n> 错误信息:%v\n\n> 串聊已超过最大文本限制，对话已重置，请重新发起。", err))
				if err != nil {
					logger.Warning(fmt.Errorf("send message error: %v", err))
					return err
				}
			} else {
				_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), fmt.Sprintf("[Wrong] 请求 OpenAI 失败了\n\n> 错误信息:%v", err))
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
			if public.JudgeSensitiveWord(reply) {
				reply = public.SolveSensitiveWord(reply)
			}
			// 回复@我的用户
			_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), FormatMarkdown(reply))
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

// FormatTimeDuation 格式化时间
// 主要提示单聊/群聊切换时多久后恢复默认聊天模式
func FormatTimeDuation(duration time.Duration) string {
	minutes := int64(duration.Minutes())
	seconds := int64(duration.Seconds()) - minutes*60
	timeoutStr := ""
	if seconds == 0 {
		timeoutStr = fmt.Sprintf("%d分钟", minutes)
	} else {
		timeoutStr = fmt.Sprintf("%d分%d秒", minutes, seconds)
	}
	return timeoutStr
}

// FormatMarkdown 格式化Markdown
// 主要修复ChatGPT返回多行代码块，钉钉会将代码块中的#当作Markdown语法里的标题来处理，进行转义；如果Markdown格式内存在html，将Markdown中的html标签转义
// 代码块缩进问题暂无法解决，因不管是四个空格，还是Tab，在钉钉上均会顶格显示，建议复制代码后用IDE进行代码格式化，针对缩进严格的语言，例如Python，不确定的建议手机端查看下代码块的缩进
func FormatMarkdown(md string) string {
	lines := strings.Split(md, "\n")
	codeblock := false
	existHtml := strings.Contains(md, "<")

	for i, line := range lines {
		if strings.HasPrefix(line, "```") {
			codeblock = !codeblock
		}
		if codeblock {
			lines[i] = strings.ReplaceAll(line, "#", "\\#")
		} else if existHtml {
			lines[i] = html.EscapeString(line)
		}
	}

	return strings.Join(lines, "\n")
}

// CheckRequestTimes 分析处理请求逻辑
// 主要提供单日请求限额的功能
func CheckRequestTimes(rmsg *dingbot.ReceiveMsg) bool {
	if public.Config.MaxRequest == 0 {
		return true
	}
	count := public.UserService.GetUseRequestCount(rmsg.GetSenderIdentifier())
	// 用户是管理员或VIP用户，不判断访问次数是否超过限制
	if public.JudgeAdminUsers(rmsg.SenderStaffId) || public.JudgeVipUsers(rmsg.SenderStaffId) {
		return true
	} else {
		// 用户不是管理员和VIP用户，判断访问次数是否超过限制
		if count >= public.Config.MaxRequest {
			logger.Info(fmt.Sprintf("亲爱的: %s，您今日请求次数已达上限，请明天再来，交互发问资源有限，请务必斟酌您的问题，给您带来不便，敬请谅解!", rmsg.SenderNick))
			_, err := rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), fmt.Sprintf("[Staple] **一个好的问题，胜过十个好的答案！** \n\n亲爱的%s:\n\n您今日请求次数已达上限，请明天再来，交互发问资源有限，请务必斟酌您的问题，给您带来不便，敬请谅解！\n\n如有需要，可联系管理员升级为VIP用户。", rmsg.SenderNick))
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
			return false
		}
	}
	// 访问次数未超过限制，将计数加1
	public.UserService.SetUseRequestCount(rmsg.GetSenderIdentifier(), count+1)
	return true
}
