package process

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/eryajf/chatgpt-dingtalk/pkg/db"
	"github.com/eryajf/chatgpt-dingtalk/pkg/dingbot"
	"github.com/eryajf/chatgpt-dingtalk/pkg/llm"
	"github.com/eryajf/chatgpt-dingtalk/pkg/logger"
	"github.com/eryajf/chatgpt-dingtalk/public"
)

// DoStream 使用流式输出执行处理请求
func DoStream(mode string, rmsg *dingbot.ReceiveMsg) error {
	// 先把模式注入
	public.UserService.SetUserMode(rmsg.GetSenderIdentifier(), mode)

	switch mode {
	case "单聊":
		return doSingleChatStream(rmsg)
	case "串聊":
		return doContextChatStream(rmsg)
	default:
		return nil
	}
}

// doSingleChatStream 单聊流式处理
func doSingleChatStream(rmsg *dingbot.ReceiveMsg) error {
	// 保存问题到数据库
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

	// 获取流式内容
	contentCh, cleanup, err := llm.SingleQaStream(rmsg.Text.Content, rmsg.GetSenderIdentifier())
	if err != nil {
		logger.Info(fmt.Errorf("gpt request error: %v", err))
		if strings.Contains(fmt.Sprintf("%v", err), "maximum question length exceeded") {
			public.UserService.ClearUserSessionContext(rmsg.GetSenderIdentifier())
			_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), fmt.Sprintf("[Wrong] 请求 OpenAI 失败了\n\n> 错误信息:%v\n\n> 已超过最大文本限制，请缩短提问文字的字数。", err))
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		} else {
			_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), fmt.Sprintf("[Wrong] 请求 OpenAI 失败了\n\n> 错误信息:%v", err))
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		}
		return err
	}
	defer cleanup()

	// 使用简化版本:直接累积内容后一次性回复
	fullContent := ""
	for content := range contentCh {
		fullContent += content
	}

	if fullContent == "" {
		logger.Warning("get gpt result failed: empty response")
		return nil
	}

	// 格式化和处理答案
	fullContent = strings.TrimSpace(fullContent)
	fullContent = strings.Trim(fullContent, "\n")

	// 保存答案到数据库
	aObj := db.Chat{
		Username:      rmsg.SenderNick,
		Source:        rmsg.GetChatTitle(),
		ChatType:      db.A,
		ParentContent: qid,
		Content:       fullContent,
	}
	_, err = aObj.Add()
	if err != nil {
		logger.Error("往MySQL新增数据失败,错误信息：", err)
	}

	logger.Info(fmt.Sprintf("🤖 %s得到的答案: %#v", rmsg.SenderNick, fullContent))

	// 敏感词过滤
	if public.JudgeSensitiveWord(fullContent) {
		fullContent = public.SolveSensitiveWord(fullContent)
	}

	// 回复用户
	_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), FormatMarkdown(fullContent))
	if err != nil {
		logger.Warning(fmt.Errorf("send message error: %v", err))
		return err
	}

	return nil
}

// doContextChatStream 串聊流式处理
func doContextChatStream(rmsg *dingbot.ReceiveMsg) error {
	// 保存问题到数据库
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

	// 获取流式内容
	cli, contentCh, err := llm.ContextQaStream(rmsg.Text.Content, rmsg.GetSenderIdentifier())
	if err != nil {
		logger.Info(fmt.Sprintf("gpt request error: %v", err))
		if strings.Contains(fmt.Sprintf("%v", err), "maximum text length exceeded") {
			public.UserService.ClearUserSessionContext(rmsg.GetSenderIdentifier())
			_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), fmt.Sprintf("[Wrong] 请求 OpenAI 失败了\n\n> 错误信息:%v\n\n> 串聊已超过最大文本限制，对话已重置，请重新发起。", err))
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		} else {
			_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), fmt.Sprintf("[Wrong] 请求 OpenAI 失败了\n\n> 错误信息:%v", err))
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		}
		return err
	}
	defer cli.Close()

	// 使用简化版本:直接累积内容后一次性回复
	fullContent := ""
	for content := range contentCh {
		fullContent += content
	}

	if fullContent == "" {
		logger.Warning("get gpt result failed: empty response")
		return nil
	}

	// 格式化和处理答案
	fullContent = strings.TrimSpace(fullContent)
	fullContent = strings.Trim(fullContent, "\n")

	// 保存答案到数据库
	aObj := db.Chat{
		Username:      rmsg.SenderNick,
		Source:        rmsg.GetChatTitle(),
		ChatType:      db.A,
		ParentContent: qid,
		Content:       fullContent,
	}
	aid, err := aObj.Add()
	if err != nil {
		logger.Error("往MySQL新增数据失败,错误信息：", err)
	}

	// 将当前回答的ID放入缓存
	public.UserService.SetAnswerID(rmsg.SenderNick, rmsg.GetChatTitle(), aid)

	logger.Info(fmt.Sprintf("🤖 %s得到的答案: %#v", rmsg.SenderNick, fullContent))

	// 敏感词过滤
	if public.JudgeSensitiveWord(fullContent) {
		fullContent = public.SolveSensitiveWord(fullContent)
	}

	// 回复用户
	_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), FormatMarkdown(fullContent))
	if err != nil {
		logger.Warning(fmt.Errorf("send message error: %v", err))
		return err
	}

	// 保存对话上下文
	_ = cli.ChatContext.SaveConversation(rmsg.GetSenderIdentifier())

	return nil
}

// DoStreamWithCard 使用流式卡片输出执行处理请求 (需要配置卡片模板)
func DoStreamWithCard(mode string, rmsg *dingbot.ReceiveMsg, cardTemplateID string) error {
	// 先把模式注入
	public.UserService.SetUserMode(rmsg.GetSenderIdentifier(), mode)

	// 检查是否有 RobotCode，如果没有则降级为简化流式模式
	clientId := rmsg.RobotCode
	if clientId == "" {
		logger.Warning("RobotCode is empty, fallback to simple stream mode")
		return DoStream(mode, rmsg)
	}

	// 获取钉钉客户端
	dingClient := public.DingTalkClientManager.GetClientByOAuthClientID(clientId)
	if dingClient == nil {
		logger.Warning(fmt.Errorf("dingtalk client not found for robot code: %s, fallback to simple stream mode", clientId))
		return DoStream(mode, rmsg)
	}

	client, ok := dingClient.(*dingbot.DingTalkClient)
	if !ok {
		logger.Warning("invalid dingtalk client type, fallback to simple stream mode")
		return DoStream(mode, rmsg)
	}

	// 生成唯一追踪ID
	trackID := uuid.New().String()

	// 创建并投放卡片
	accessToken, err := client.GetAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	cardClient, err := dingbot.NewStreamCardClient()
	if err != nil {
		return fmt.Errorf("failed to create stream card client: %w", err)
	}

	// 构建OpenSpaceID
	var openSpaceID string
	if rmsg.ConversationType == "2" { // 群聊
		openSpaceID = fmt.Sprintf("dtv1.card//IM_GROUP.%s", rmsg.ConversationID)
		logger.Info(fmt.Sprintf("🎴 群聊模式 - OpenSpaceID: %s, RobotCode: %s", openSpaceID, rmsg.RobotCode))
	} else { // 单聊
		openSpaceID = fmt.Sprintf("dtv1.card//IM_ROBOT.%s", rmsg.SenderStaffId)
		logger.Info(fmt.Sprintf("🎴 私聊模式 - OpenSpaceID: %s, ConversationType: %s", openSpaceID, rmsg.ConversationType))
	}

	createReq := &dingbot.CreateAndDeliverCardRequest{
		CardTemplateID:   cardTemplateID,
		OutTrackID:       trackID,
		ConversationID:   rmsg.ConversationID,
		SenderStaffID:    rmsg.SenderStaffId,
		RobotCode:        rmsg.RobotCode,
		OpenSpaceID:      openSpaceID,
		ConversationType: rmsg.ConversationType,
		CardData: map[string]string{
			"content": "",
		},
	}

	if err := cardClient.CreateAndDeliverCard(accessToken, createReq); err != nil {
		logger.Warning(fmt.Errorf("failed to create card: %v", err))
		// 卡片创建失败,降级为普通消息
		return DoStream(mode, rmsg)
	}

	// 发送初始状态
	initialContent := fmt.Sprintf("**%s**\n\n%s", rmsg.Text.Content, "稍等，让我想一想……")
	if err := client.UpdateAIStreamCard(trackID, initialContent, false); err != nil {
		logger.Warning(fmt.Errorf("failed to update initial card: %v", err))
	}

	// 获取流式内容
	var contentCh <-chan string
	var cli *llm.Client
	if mode == "单聊" {
		var cleanup func()
		contentCh, cleanup, err = llm.SingleQaStream(rmsg.Text.Content, rmsg.GetSenderIdentifier())
		defer cleanup()
	} else {
		cli, contentCh, err = llm.ContextQaStream(rmsg.Text.Content, rmsg.GetSenderIdentifier())
		defer cli.Close()
	}

	if err != nil {
		errorMsg := fmt.Sprintf("**%s**\n\n出错了: %v", rmsg.Text.Content, err)
		if err := client.UpdateAIStreamCard(trackID, errorMsg, true); err != nil {
			logger.Warning(fmt.Errorf("failed to update error card: %v", err))
		}
		return err
	}

	// 实时流式更新卡片内容
	questionHeader := fmt.Sprintf("**%s**\n\n", rmsg.Text.Content)
	fullContent := questionHeader

	// 使用缓冲机制避免更新过于频繁
	updateBuffer := ""
	lastUpdateTime := time.Now()
	minUpdateInterval := 300 * time.Millisecond // 最小更新间隔300ms

	for {
		content, ok := <-contentCh
		if !ok {
			// 流结束,发送最后的更新(如果有未发送的缓冲内容)
			if updateBuffer != "" {
				fullContent += updateBuffer
				if err := client.UpdateAIStreamCard(trackID, fullContent, true); err != nil {
					logger.Error(fmt.Errorf("failed to finalize card: %v", err))
				}
			} else {
				// 标记为完成
				if err := client.UpdateAIStreamCard(trackID, fullContent, true); err != nil {
					logger.Error(fmt.Errorf("failed to finalize card: %v", err))
				}
			}

			// 保存到数据库并处理后续逻辑
			saveStreamResult(mode, rmsg, fullContent[len(questionHeader):], cli)
			return nil
		}

		// 累积接收到的内容到缓冲区
		updateBuffer += content

		// 检查是否应该更新(距离上次更新超过最小间隔)
		if time.Since(lastUpdateTime) >= minUpdateInterval {
			fullContent += updateBuffer
			updateBuffer = ""

			// 立即更新卡片
			if err := client.UpdateAIStreamCard(trackID, fullContent, false); err != nil {
				logger.Warning(fmt.Errorf("failed to update card: %v", err))
			}

			lastUpdateTime = time.Now()
		}
	}
}

// saveStreamResult 保存流式结果到数据库
func saveStreamResult(mode string, rmsg *dingbot.ReceiveMsg, answer string, cli *llm.Client) {
	answer = strings.TrimSpace(answer)
	answer = strings.Trim(answer, "\n")

	if mode == "单聊" {
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

		aObj := db.Chat{
			Username:      rmsg.SenderNick,
			Source:        rmsg.GetChatTitle(),
			ChatType:      db.A,
			ParentContent: qid,
			Content:       answer,
		}
		_, err = aObj.Add()
		if err != nil {
			logger.Error("往MySQL新增数据失败,错误信息：", err)
		}
	} else { // 串聊
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

		aObj := db.Chat{
			Username:      rmsg.SenderNick,
			Source:        rmsg.GetChatTitle(),
			ChatType:      db.A,
			ParentContent: qid,
			Content:       answer,
		}
		aid, err := aObj.Add()
		if err != nil {
			logger.Error("往MySQL新增数据失败,错误信息：", err)
		}

		public.UserService.SetAnswerID(rmsg.SenderNick, rmsg.GetChatTitle(), aid)

		if cli != nil {
			_ = cli.ChatContext.SaveConversation(rmsg.GetSenderIdentifier())
		}
	}

	logger.Info(fmt.Sprintf("🤖 %s得到的答案: %#v", rmsg.SenderNick, answer))
}
