package public

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

func InitAiCli() *resty.Client {
	if Config.HttpProxy != "" {
		return resty.New().SetTimeout(10*time.Second).SetHeader("Authorization", fmt.Sprintf("Bearer %s", Config.ApiKey)).SetProxy(Config.HttpProxy).SetRetryCount(3).SetRetryWaitTime(2 * time.Second)
	}
	return resty.New().SetTimeout(10*time.Second).SetHeader("Authorization", fmt.Sprintf("Bearer %s", Config.ApiKey)).SetRetryCount(3).SetRetryWaitTime(2 * time.Second)
}

type Bill struct {
	Object     string      `json:"object"`
	DailyCosts []DailyCost `json:"daily_costs"`
	TotalUsage float64     `json:"total_usage"`
}

type DailyCost struct {
	Timestamp float64    `json:"timestamp"`
	LineItems []LineItem `json:"line_items"`
}

type LineItem struct {
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

// GetBalance 获取账号余额
func GetBalance() (string, error) {
	var data Bill
	path := "/v1/dashboard/billing/usage"
	var url string = "https://api.openai.com" + path
	if Config.BaseURL != "" {
		url = Config.BaseURL + path
	}
	d, _ := time.ParseDuration("-24h")
	resp, err := InitAiCli().R().SetQueryParams(map[string]string{
		"start_date": time.Now().Add(d * 90).Format("2006-01-02"),
		"end_date":   time.Now().Format("2006-01-02"),
	}).Get(url)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return "", err
	}
	sub, err := GetSub()
	if err != nil {
		return "", err
	}
	expireDate := time.Unix(sub.AccessUntil, 0).Format("2006-01-02 15:04:05")
	used := data.TotalUsage / 100
	totalAvailable := sub.HardLimitUsd - used
	msg := fmt.Sprintf("💵 已用: 💲%v\n💵 剩余: 💲%v\n🕰 到期时间: %v", fmt.Sprintf("%.2f", used), fmt.Sprintf("%.2f", totalAvailable), expireDate)
	// 放入缓存
	UserService.SetUserMode("system_balance", msg)
	return msg, nil
}

type Subscription struct {
	Object             string      `json:"object"`
	HasPaymentMethod   bool        `json:"has_payment_method"`
	Canceled           bool        `json:"canceled"`
	CanceledAt         interface{} `json:"canceled_at"`
	Delinquent         interface{} `json:"delinquent"`
	AccessUntil        int64       `json:"access_until"`
	SoftLimit          int64       `json:"soft_limit"`
	HardLimit          int64       `json:"hard_limit"`
	SystemHardLimit    int64       `json:"system_hard_limit"`
	SoftLimitUsd       float64     `json:"soft_limit_usd"`
	HardLimitUsd       float64     `json:"hard_limit_usd"`
	SystemHardLimitUsd float64     `json:"system_hard_limit_usd"`
	Plan               Plan        `json:"plan"`
	AccountName        string      `json:"account_name"`
	PoNumber           interface{} `json:"po_number"`
	BillingEmail       interface{} `json:"billing_email"`
	TaxIDS             interface{} `json:"tax_ids"`
	BillingAddress     interface{} `json:"billing_address"`
	BusinessAddress    interface{} `json:"business_address"`
}

type Plan struct {
	Title string `json:"title"`
	ID    string `json:"id"`
}

func GetSub() (Subscription, error) {
	var data Subscription
	path := "/v1/dashboard/billing/subscription"
	var url string = "https://api.openai.com" + path
	if Config.BaseURL != "" {
		url = Config.BaseURL + path
	}
	resp, err := InitAiCli().R().Get(url)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}
