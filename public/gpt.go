package public

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

func InitAiCli(apiKey string) *resty.Client {
	if Config.HttpProxy != "" {
		return resty.New().SetTimeout(30*time.Second).SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey)).SetProxy(Config.HttpProxy).SetRetryCount(3).SetRetryWaitTime(5 * time.Second)
	}
	return resty.New().SetTimeout(30*time.Second).SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey)).SetRetryCount(3).SetRetryWaitTime(5 * time.Second)
}

type Billing struct {
	Object         string  `json:"object"`
	TotalGranted   float64 `json:"total_granted"`
	TotalUsed      float64 `json:"total_used"`
	TotalAvailable float64 `json:"total_available"`
	Grants         struct {
		Object string `json:"object"`
		Data   []struct {
			Object      string  `json:"object"`
			ID          string  `json:"id"`
			GrantAmount float64 `json:"grant_amount"`
			UsedAmount  float64 `json:"used_amount"`
			EffectiveAt float64 `json:"effective_at"`
			ExpiresAt   float64 `json:"expires_at"`
		} `json:"data"`
	} `json:"grants"`
}

type ErrorResp struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	}
}

func GetBalance(apiKey string) (Billing, error) {
	var data Billing
	url := "https://api.openai.com/dashboard/billing/credit_grants"
	resp, err := InitAiCli(apiKey).R().Get(url)
	if err != nil {
		return data, err
	}

	if resp.StatusCode() != 200 {
		var errorResp ErrorResp
		err = json.Unmarshal(resp.Body(), &errorResp)
		if err != nil {
			return data, err
		}
		return data, fmt.Errorf("error: %v", errorResp.Error.Message)
	}

	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}
