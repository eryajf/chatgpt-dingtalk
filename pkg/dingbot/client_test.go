package dingbot

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/eryajf/chatgpt-dingtalk/config"
)

func TestUploadMedia_Pass_WithValidConfig(t *testing.T) {
	// 设置了钉钉 ClientID 和 ClientSecret 的环境变量才执行以下测试，用于快速验证钉钉图片上传能力
	clientId, clientSecret := os.Getenv("DINGTALK_CLIENT_ID_FOR_TEST"), os.Getenv("DINGTALK_CLIENT_SECRET_FOR_TEST")
	if len(clientId) <= 0 || len(clientSecret) <= 0 {
		return
	}
	credentials := []config.Credential{
		config.Credential{
			ClientID:     clientId,
			ClientSecret: clientSecret,
		},
	}
	client := NewDingTalkClientManager(&config.Configuration{Credentials: credentials}).GetClientByOAuthClientID(clientId)
	var imageContent []byte
	{
		// 生成一张用于测试的图片
		img := image.NewRGBA(image.Rect(0, 0, 200, 100))
		blue := color.RGBA{0, 0, 255, 255}
		for x := 0; x < img.Bounds().Dx(); x++ {
			for y := 0; y < img.Bounds().Dy(); y++ {
				img.Set(x, y, blue)
			}
		}
		buf := new(bytes.Buffer)
		err := png.Encode(buf, img)
		if err != nil {
			return
		}
		// get the byte array from the buffer
		imageContent = buf.Bytes()
	}
	result, err := client.UploadMedia(imageContent, "filename.png", "image", "image/png")
	if err != nil {
		t.Errorf("upload media failed, err=%s", err.Error())
		return
	}
	if result.MediaID == "" {
		t.Errorf("upload media failed, empty media id")
		return
	}
}
