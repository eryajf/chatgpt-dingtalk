package dingbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	url2 "net/url"
	"sync"
	"time"

	"github.com/eryajf/chatgpt-dingtalk/config"
)

// OpenAPI doc: https://open.dingtalk.com/document/isvapp/upload-media-files
const (
	MediaTypeImage string = "image"
	MediaTypeVoice string = "voice"
	MediaTypeVideo string = "video"
	MediaTypeFile  string = "file"
)
const (
	MimeTypeImagePng string = "image/png"
)

type MediaUploadResult struct {
	ErrorCode    int64  `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
	MediaID      string `json:"media_id"`
	CreatedAt    int64  `json:"created_at"`
	Type         string `json:"type"`
}

type OAuthTokenResult struct {
	ErrorCode    int    `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type DingTalkClientInterface interface {
	GetAccessToken() (string, error)
	UploadMedia(content []byte, filename, mediaType, mimeType string) (*MediaUploadResult, error)
}

type DingTalkClientManagerInterface interface {
	GetClientByOAuthClientID(clientId string) DingTalkClientInterface
}

type DingTalkClient struct {
	Credential  config.Credential
	AccessToken string
	expireAt    int64
	mutex       sync.Mutex
}

type DingTalkClientManager struct {
	Credentials []config.Credential
	Clients     map[string]*DingTalkClient
	mutex       sync.Mutex
}

func NewDingTalkClient(credential config.Credential) *DingTalkClient {
	return &DingTalkClient{
		Credential: credential,
	}
}

func NewDingTalkClientManager(conf *config.Configuration) *DingTalkClientManager {
	clients := make(map[string]*DingTalkClient)

	if conf != nil && conf.Credentials != nil {
		for _, credential := range conf.Credentials {
			clients[credential.ClientID] = NewDingTalkClient(credential)
		}
	}
	return &DingTalkClientManager{
		Credentials: conf.Credentials,
		Clients:     clients,
	}
}

func (m *DingTalkClientManager) GetClientByOAuthClientID(clientId string) DingTalkClientInterface {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if client, ok := m.Clients[clientId]; ok {
		return client
	}
	return nil
}

func (c *DingTalkClient) GetAccessToken() (string, error) {
	accessToken := ""
	{
		// 先查询缓存
		c.mutex.Lock()
		now := time.Now().Unix()
		if c.expireAt > 0 && c.AccessToken != "" && (now+60) < c.expireAt {
			// 预留一分钟有效期避免在Token过期的临界点调用接口出现401错误
			accessToken = c.AccessToken
		}
		c.mutex.Unlock()
	}
	if accessToken != "" {
		return accessToken, nil
	}

	tokenResult, err := c.getAccessTokenFromDingTalk()
	if err != nil {
		return "", err
	}

	{
		// 更新缓存
		c.mutex.Lock()
		c.AccessToken = tokenResult.AccessToken
		c.expireAt = time.Now().Unix() + int64(tokenResult.ExpiresIn)
		c.mutex.Unlock()
	}
	return tokenResult.AccessToken, nil
}

func (c *DingTalkClient) UploadMedia(content []byte, filename, mediaType, mimeType string) (*MediaUploadResult, error) {
	// OpenAPI doc: https://open.dingtalk.com/document/isvapp/upload-media-files
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}
	if len(accessToken) == 0 {
		return nil, errors.New("empty access token")
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("media", filename)
	if err != nil {
		return nil, err
	}
	_, err = part.Write(content)
	if err != nil {
		return nil, err
	}
	if err = writer.WriteField("type", mediaType); err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Create a new HTTP request to upload the media file
	url := fmt.Sprintf("https://oapi.dingtalk.com/media/upload?access_token=%s", url2.QueryEscape(accessToken))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the HTTP request and parse the response
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse the response body as JSON and extract the media ID
	media := &MediaUploadResult{}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bodyBytes, media); err != nil {
		return nil, err
	}
	if media.ErrorCode != 0 {
		return nil, errors.New(media.ErrorMessage)
	}
	return media, nil
}

func (c *DingTalkClient) getAccessTokenFromDingTalk() (*OAuthTokenResult, error) {
	// OpenAPI doc: https://open.dingtalk.com/document/orgapp/obtain-orgapp-token
	apiUrl := "https://oapi.dingtalk.com/gettoken"
	queryParams := url2.Values{}
	queryParams.Add("appkey", c.Credential.ClientID)
	queryParams.Add("appsecret", c.Credential.ClientSecret)

	// Create a new HTTP request to get the AccessToken
	req, err := http.NewRequest("GET", apiUrl+"?"+queryParams.Encode(), nil)
	if err != nil {
		return nil, err
	}

	// Send the HTTP request and parse the response body as JSON
	client := http.Client{
		Timeout: time.Second * 60,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	tokenResult := &OAuthTokenResult{}
	err = json.Unmarshal(body, tokenResult)
	if err != nil {
		return nil, err
	}
	if tokenResult.ErrorCode != 0 {
		return nil, errors.New(tokenResult.ErrorMessage)
	}
	return tokenResult, nil
}
