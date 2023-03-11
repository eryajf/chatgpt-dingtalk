package public

import (
	"fmt"
	"log"
	"time"

	"github.com/eryajf/chatgpt-dingtalk/public/logger"
)

type ApiKeyInfo struct {
	ApiKey         string
	TotalUsed      float64
	TotalAvailable float64
	T1             time.Time
	T2             time.Time
}

type ApiKeyInfoList struct {
	AllTotalUsed      float64
	AllTotalAvailable float64
	NowUsedIndex      int
	ApiKeyInfoList    []ApiKeyInfo
}

func getNewApiKeyInfo(ApiKey string) (ApiKeyInfo, error) {

	rst, err := GetBalance(ApiKey)
	if err != nil {
		logger.Warning(fmt.Errorf("get balance error: %v", err))
		return ApiKeyInfo{}, err
	}

	t1 := time.Unix(int64(rst.Grants.Data[0].EffectiveAt), 0)
	t2 := time.Unix(int64(rst.Grants.Data[0].ExpiresAt), 0)

	return ApiKeyInfo{
		ApiKey:         ApiKey,
		TotalUsed:      rst.TotalUsed,
		TotalAvailable: rst.TotalAvailable,
		T1:             t1,
		T2:             t2,
	}, nil
}

func InitApiKeyInfo() *ApiKeyInfoList {
	apiKeyList := &ApiKeyInfoList{}

	for _, v := range Config.ApiKey {
		apiKeyInfo, err := getNewApiKeyInfo(v)
		if err != nil {
			continue
		}
		apiKeyList.ApiKeyInfoList = append(apiKeyList.ApiKeyInfoList, apiKeyInfo)

		apiKeyList.AllTotalUsed += apiKeyInfo.TotalUsed
		apiKeyList.AllTotalAvailable += apiKeyInfo.TotalAvailable

		msg := fmt.Sprintf("💵 当前key: %v**** ", v[0:16])
		msg += fmt.Sprintf("💵 当前已用: %v ", fmt.Sprintf("%.8f", apiKeyInfo.TotalUsed))
		msg += fmt.Sprintf("💵 当前剩余: %v ", fmt.Sprintf("%.8f", apiKeyInfo.TotalAvailable))
		msg += fmt.Sprintf("💵 当前有效时间: 从 %v 到 %v", apiKeyInfo.T1.Format("2006-01-02"), apiKeyInfo.T2.Format("2006-01-02"))
		log.Println(msg)
	}

	msg := fmt.Sprintf("💵 配置数量: %v ", len(Config.ApiKey))
	msg += fmt.Sprintf("💵 有效数量: %v ", len(apiKeyList.ApiKeyInfoList))
	msg += fmt.Sprintf("💵 总共已用: %v ", fmt.Sprintf("%.8f", apiKeyList.AllTotalUsed))
	msg += fmt.Sprintf("💵 总共剩余: %v ", fmt.Sprintf("%.8f", apiKeyList.AllTotalAvailable))
	log.Println(msg)

	return apiKeyList
}

func (apiKeyList *ApiKeyInfoList) GetApiKeyInfoString() string {
	var msg string

	msg += fmt.Sprintf("💵 全部已用: %v\n", fmt.Sprintf("%.8f", apiKeyList.AllTotalUsed))
	msg += fmt.Sprintf("💵 全部剩余: %v\n\n", fmt.Sprintf("%.8f", apiKeyList.AllTotalAvailable))
	msg += fmt.Sprintf("💵 当前索引: %v\n", apiKeyList.NowUsedIndex)

	// 获取余额时, 立即更新
	nowApiKey := &apiKeyList.ApiKeyInfoList[apiKeyList.NowUsedIndex]
	apiKeyInfo, err := getNewApiKeyInfo(nowApiKey.ApiKey)
	if err != nil {
		logger.Warning(fmt.Errorf("get balance error: %v", err))
	}
	nowApiKey = &apiKeyInfo

	msg += fmt.Sprintf("💵 当前已用: %v\n", fmt.Sprintf("%.8f", nowApiKey.TotalUsed))
	msg += fmt.Sprintf("💵 当前剩余: %v\n", fmt.Sprintf("%.8f", nowApiKey.TotalAvailable))
	msg += fmt.Sprintf("💵 当前周期:  %v == %v\n", nowApiKey.T1.Format("2006-01-02"), nowApiKey.T2.Format("2006-01-02"))

	return msg
}

func (apiKeyList *ApiKeyInfoList) GetApiKey(next bool) (string, error) {
	nowApiKey := &apiKeyList.ApiKeyInfoList[apiKeyList.NowUsedIndex]

	if next {
		// TODO: 当前移动到最后不会重复尝试, 下次再尝试?
		apiKeyList.NowUsedIndex++
	}

	if apiKeyList.NowUsedIndex > len(apiKeyList.ApiKeyInfoList)-1 {
		return "", fmt.Errorf("no api key available %v [0/%v]", apiKeyList.NowUsedIndex, len(apiKeyList.ApiKeyInfoList))
	}

	return nowApiKey.ApiKey, nil
}
