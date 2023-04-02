package public

import (
	"io/ioutil"
	"os"
	"strings"
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
	if len(Config.AllowGroups) == 0 {
		return false
	}
	for _, v := range Config.AdminUsers {
		if v == s {
			return true
		}
	}
	return false
}
