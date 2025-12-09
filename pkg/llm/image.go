package llm

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"strings"
	"time"

	"golang.org/x/image/webp"

	openai "github.com/sashabaranov/go-openai"

	"github.com/eryajf/chatgpt-dingtalk/pkg/dingbot"
	"github.com/eryajf/chatgpt-dingtalk/public"
)

func getImageTypeFromBase64(base64Str string) string {
	switch {
	case strings.HasPrefix(base64Str, "/9j/"):
		return "JPEG"
	case strings.HasPrefix(base64Str, "iVBOR"):
		return "PNG"
	case strings.HasPrefix(base64Str, "R0lG"):
		return "GIF"
	case strings.HasPrefix(base64Str, "UklG"):
		return "WebP"
	default:
		return "Unknown"
	}
}

func (c *Client) GenerateImage(ctx context.Context, prompt string) (string, error) {
	model := public.Config.Model
	imageModel := public.Config.ImageModel
	if isModelSupportedChatCompletions(model) {
		req := openai.ImageRequest{
			Prompt:         prompt,
			Model:          imageModel,
			Size:           openai.CreateImageSize1024x1024,
			ResponseFormat: openai.CreateImageResponseFormatB64JSON,
			N:              1,
			User:           c.userId,
		}
		respBase64, err := c.client.CreateImage(c.ctx, req)
		if err != nil {
			return "", err
		}
		imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
		if err != nil {
			return "", err
		}

		r := bytes.NewReader(imgBytes)

		imgType := getImageTypeFromBase64(respBase64.Data[0].B64JSON)
		var imgData image.Image
		var imgErr error
		if imgType == "WebP" {
			imgData, imgErr = webp.Decode(r)
		} else {
			imgData, _, imgErr = image.Decode(r)
		}
		if imgErr != nil {
			return "", imgErr
		}

		imageName := time.Now().Format("20060102-150405") + ".png"
		clientId, _ := ctx.Value(public.DingTalkClientIdKeyName).(string)
		client := public.DingTalkClientManager.GetClientByOAuthClientID(clientId)
		mediaResult, uploadErr := &dingbot.MediaUploadResult{}, errors.New(fmt.Sprintf("unknown clientId: %s", clientId))
		if client != nil {
			mediaResult, uploadErr = client.UploadMedia(imgBytes, imageName, dingbot.MediaTypeImage, dingbot.MimeTypeImagePng)
		}

		err = os.MkdirAll("data/images", 0755)
		if err != nil {
			return "", err
		}
		file, err := os.Create("data/images/" + imageName)
		if err != nil {
			return "", err
		}
		defer file.Close()

		if err := png.Encode(file, imgData); err != nil {
			return "", err
		}
		if uploadErr == nil {
			return mediaResult.MediaID, nil
		} else {
			return public.Config.ServiceURL + "/images/" + imageName, nil
		}
	}
	return "", nil
}
