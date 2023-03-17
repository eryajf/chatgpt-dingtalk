package process

import (
	"fmt"
	"strings"
	"time"

	"github.com/eryajf/chatgpt-dingtalk/public"
	"github.com/eryajf/chatgpt-dingtalk/public/logger"
	"github.com/solywsh/chatgpt"
)

// ProcessRequest 分析处理请求逻辑
func ProcessRequest(rmsg public.ReceiveMsg) error {
	if CheckRequest(rmsg) {
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
	}
	return nil
}

// 执行处理请求
func Do(mode string, rmsg public.ReceiveMsg) error {
	// 先把模式注入
	public.UserService.SetUserMode(rmsg.SenderStaffId, mode)
	switch mode {
	case "单聊":
		reply, err := chatgpt.SingleQa(rmsg.Text.Content, rmsg.SenderStaffId)
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
		cli, reply, err := chatgpt.ContextQa(rmsg.Text.Content, rmsg.SenderStaffId)
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

// ProcessRequest 分析处理请求逻辑
func CheckRequest(rmsg public.ReceiveMsg) bool {
	if public.Config.MaxRequest == 0 {
		return true
	}
	count := public.UserService.GetUseRequestCount(rmsg.SenderStaffId)
	// 判断访问次数是否超过限制
	if count >= public.Config.MaxRequest {
		logger.Info(fmt.Sprintf("亲爱的: %s，您今日请求次数已达上限，请明天再来，交互发问资源有限，请务必斟酌您的问题，给您带来不便，敬请谅解!", rmsg.SenderNick))
		_, err := rmsg.ReplyText(fmt.Sprintf("一个好的问题，胜过十个好的答案！\n亲爱的: %s，您今日请求次数已达上限，请明天再来，交互发问资源有限，请务必斟酌您的问题，给您带来不便，敬请谅解!", rmsg.SenderNick), rmsg.SenderStaffId)
		if err != nil {
			logger.Warning(fmt.Errorf("send message error: %v", err))
		}
		return false
	}
	// 访问次数未超过限制，将计数加1
	public.UserService.SetUseRequestCount(rmsg.SenderStaffId, count+1)
	return true
}
