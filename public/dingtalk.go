package public

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// 接收的消息体
type ReceiveMsg struct {
	ConversationID string `json:"conversationId"`
	AtUsers        []struct {
		DingtalkID string `json:"dingtalkId"`
	} `json:"atUsers"`
	ChatbotUserID             string `json:"chatbotUserId"`
	MsgID                     string `json:"msgId"`
	SenderNick                string `json:"senderNick"`
	IsAdmin                   bool   `json:"isAdmin"`
	SenderStaffId             string `json:"senderStaffId"`
	SessionWebhookExpiredTime int64  `json:"sessionWebhookExpiredTime"`
	CreateAt                  int64  `json:"createAt"`
	ConversationType          string `json:"conversationType"`
	SenderID                  string `json:"senderId"`
	ConversationTitle         string `json:"conversationTitle"`
	IsInAtList                bool   `json:"isInAtList"`
	SessionWebhook            string `json:"sessionWebhook"`
	Text                      Text   `json:"text"`
	RobotCode                 string `json:"robotCode"`
	Msgtype                   string `json:"msgtype"`
}


// 发送的消息体
type SendMsg struct {
	Text    Text   `json:"text"`
	Msgtype string `json:"msgtype"`
	At 		At `json:"at"`
}

// 消息内容
type Text struct {
	Content string `json:"content"`
}

// at 内容
type At struct {
	AtUserIds []string `json:"atUserIds"`
}

// 发消息给钉钉
func (r ReceiveMsg) ReplyText(msg string, staffId string) (statuscode int, err error) {
	// 定义消息
	msgtmp := &SendMsg{Text: Text{Content: msg}, Msgtype: "text", At: At{AtUserIds: []string{staffId}}}
	data, err := json.Marshal(msgtmp)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", r.SessionWebhook, bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}
