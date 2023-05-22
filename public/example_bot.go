package public

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/payload"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/utils"
)

/**
 * @Author linya.jj
 * @Date 2023/3/22 18:30
 */

func OnBotCallback(ctx context.Context, df *payload.DataFrame) (*payload.DataFrameResponse, error) {
	frameResp := &payload.DataFrameResponse{
		Code: 200,
		Headers: payload.DataFrameHeader{
			payload.DataFrameHeaderKContentType: payload.DataFrameContentTypeKJson,
			payload.DataFrameHeaderKMessageId:   df.GetMessageId(),
		},
		Message: "ok",
		Data:    "",
	}

	return frameResp, nil
}

func OnChatReceive(ctx context.Context, data *chatbot.BotCallbackDataModel) error {
	requestBody := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": fmt.Sprintf("msg received: [%s]", data.Text.Content),
		},
	}

	requestJsonBody, _ := json.Marshal(requestBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, data.SessionWebhook, bytes.NewReader(requestJsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")

	httpClient := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   5 * time.Second, //设置超时，包含connection时间、任意重定向时间、读取response body时间
	}

	_, err = httpClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func RunBotListener(clientId, clientSecret string) {
	logger.SetLogger(logger.NewStdTestLogger())

	cli := client.NewStreamClient(
		client.WithAppCredential(client.NewAppCredentialConfig(clientId, clientSecret)),
		client.WithUserAgent(client.NewDingtalkGoSDKUserAgent()),
		client.WithSubscription(utils.SubscriptionTypeKCallback, payload.BotMessageCallbackTopic, chatbot.NewDefaultChatBotFrameHandler(OnChatReceive).OnEventReceived),
	)

	err := cli.Start(context.Background())
	if err != nil {
		panic(err)
	}

	defer cli.Close()

	select {}
}
