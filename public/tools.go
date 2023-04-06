package public

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// 将内容写入到文件，如果文件名带路径，则会判断路径是否存在，不存在则创建
func WriteToFile(path string, data []byte) error {
	tmp := strings.Split(path, "/")
	if len(tmp) > 0 {
		tmp = tmp[:len(tmp)-1]
	}

	err := os.MkdirAll(strings.Join(tmp, "/"), os.ModePerm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0755)
	if err != nil {
		return err
	}
	return nil
}

// JudgeGroup 判断群聊名称是否在白名单
func JudgeGroup(s string) bool {
	if len(Config.AllowGroups) == 0 {
		return true
	}
	for _, v := range Config.AllowGroups {
		if v == s {
			return true
		}
	}
	return false
}

// JudgeUsers 判断用户名称是否在白名单
func JudgeUsers(s string) bool {
	if len(Config.AllowUsers) == 0 {
		return true
	}
	for _, v := range Config.AllowUsers {
		if v == s {
			return true
		}
	}
	return false
}

// JudgeAdminUsers 判断用户是否为系统管理员
func JudgeAdminUsers(s string) bool {
	// 如果secret或者用户的userid都为空的话，那么默认没有管理员
	if len(Config.AppSecrets) == 0 || s == "" {
		return false
	}
	// 如果没有指定，则没有人是管理员
	if len(Config.AdminUsers) == 0 {
		return false
	}
	for _, v := range Config.AdminUsers {
		if v == s {
			return true
		}
	}
	return false
}

func GetReadTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func CheckRequest(ts, sg string) bool {
	appSecrets := Config.AppSecrets
	// 如果没有指定或者outgoing类型机器人下使用，则默认不做校验
	if len(appSecrets) == 0 || sg == "" {
		return true
	}
	// 校验appSecret
	for _, secret := range appSecrets {
		stringToSign := fmt.Sprintf("%s\n%s", ts, secret)
		mac := hmac.New(sha256.New, []byte(secret))
		_, _ = mac.Write([]byte(stringToSign))
		if base64.StdEncoding.EncodeToString(mac.Sum(nil)) == sg {
			return true
		}
	}
	return false
}
