package chatgpt

import openai "github.com/sashabaranov/go-openai"

var ModelsSupportChatCompletions = []string{
	openai.GPT432K0613,
	openai.GPT432K0314,
	openai.GPT432K,
	openai.GPT40613,
	openai.GPT40314,
	openai.GPT4TurboPreview,
	openai.GPT4VisionPreview,
	openai.GPT4,
	openai.GPT4oMini,
	openai.GPT3Dot5Turbo1106,
	openai.GPT3Dot5Turbo0613,
	openai.GPT3Dot5Turbo0301,
	openai.GPT3Dot5Turbo16K,
	openai.GPT3Dot5Turbo16K0613,
	openai.GPT3Dot5Turbo,
}

func isModelSupportedChatCompletions(model string) bool {
	for _, m := range ModelsSupportChatCompletions {
		if m == model {
			return true
		}
	}
	return false
}
