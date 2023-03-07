package public

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

func InitAiCli() *resty.Client {
	if Config.HttpProxy != "" {
		return resty.New().SetTimeout(30*time.Second).SetHeader("Authorization", fmt.Sprintf("Bearer %s", Config.ApiKey)).SetProxy(Config.HttpProxy).SetRetryCount(3).SetRetryWaitTime(5 * time.Second)
	}
	return resty.New().SetTimeout(30*time.Second).SetHeader("Authorization", fmt.Sprintf("Bearer %s", Config.ApiKey)).SetRetryCount(3).SetRetryWaitTime(5 * time.Second)
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

func GetBalance() (Billing, error) {
	var data Billing
	url := "https://api.openai.com/dashboard/billing/credit_grants"
	resp, err := InitAiCli().R().Get(url)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return data, err
	}
	t1 := time.Unix(int64(data.Grants.Data[0].EffectiveAt), 0)
	t2 := time.Unix(int64(data.Grants.Data[0].ExpiresAt), 0)
	msg := fmt.Sprintf("💵 已用: 💲%v\n💵 剩余: 💲%v\n⏳ 有效时间: 从 %v 到 %v\n", fmt.Sprintf("%.2f", data.TotalUsed), fmt.Sprintf("%.2f", data.TotalAvailable), t1.Format("2006-01-02 15:04:05"), t2.Format("2006-01-02 15:04:05"))
	// 放入缓存
	UserService.SetUserMode("system_balance", msg)
	return data, nil
}
