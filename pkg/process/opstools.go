package process

import (
	"fmt"
	"strings"

	"github.com/eryajf/chatgpt-dingtalk/pkg/db"
	"github.com/eryajf/chatgpt-dingtalk/pkg/dingbot"
	"github.com/eryajf/chatgpt-dingtalk/pkg/logger"
	"github.com/eryajf/chatgpt-dingtalk/pkg/ops"
	"github.com/eryajf/chatgpt-dingtalk/public"
)

// 一些运维方面的工具在此

// 域名信息
func DomainMsg(rmsg *dingbot.ReceiveMsg) error {
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
	domain := strings.TrimSpace(strings.Split(rmsg.Text.Content, " ")[1])
	dm, err := ops.GetDomainMsg(domain)
	if err != nil {
		return err
	}
	// 回复@我的用户
	reply := fmt.Sprintf("**创建时间:** %v\n\n**到期时间:** %v\n\n**服务商:** %v", dm.CreateDate, dm.ExpiryDate, dm.Registrar)
	aObj := db.Chat{
		Username:      rmsg.SenderNick,
		Source:        rmsg.GetChatTitle(),
		ChatType:      db.A,
		ParentContent: qid,
		Content:       reply,
	}
	_, err = aObj.Add()
	if err != nil {
		logger.Error("往MySQL新增数据失败,错误信息：", err)
	}
	logger.Info(fmt.Sprintf("🤖 %s得到的答案: %#v", rmsg.SenderNick, reply))
	_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), reply)
	if err != nil {
		logger.Error(fmt.Errorf("send message error: %v", err))
		return err
	}
	return nil
}

// 证书信息
func DomainCertMsg(rmsg *dingbot.ReceiveMsg) error {
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
	domain := strings.TrimSpace(strings.Split(rmsg.Text.Content, " ")[1])
	dm, err := ops.GetDomainCertMsg(domain)
	if err != nil {
		return err
	}
	cert := dm.PeerCertificates[0]
	// 回复@我的用户
	reply := fmt.Sprintf("**证书创建时间:** %v\n\n**证书到期时间:** %v\n\n**证书颁发机构:** %v\n\n", public.GetReadTime(cert.NotBefore), public.GetReadTime(cert.NotAfter), cert.Issuer.Organization)
	aObj := db.Chat{
		Username:      rmsg.SenderNick,
		Source:        rmsg.GetChatTitle(),
		ChatType:      db.A,
		ParentContent: qid,
		Content:       reply,
	}
	_, err = aObj.Add()
	if err != nil {
		logger.Error("往MySQL新增数据失败,错误信息：", err)
	}
	logger.Info(fmt.Sprintf("🤖 %s得到的答案: %#v", rmsg.SenderNick, reply))
	_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), reply)
	if err != nil {
		logger.Error(fmt.Errorf("send message error: %v", err))
		return err
	}
	return nil
}
