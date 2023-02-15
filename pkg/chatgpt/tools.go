package chatgpt

import (
	"fmt"
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
	fmt.Println(tmp)
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
