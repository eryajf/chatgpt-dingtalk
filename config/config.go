package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

// Configuration 项目配置
type Configuration struct {
	// gtp apikey
	ApiKey string `json:"api_key"`
	// 会话超时时间
	SessionTimeout time.Duration `json:"session_timeout"`
}

var config *Configuration
var once sync.Once

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	once.Do(func() {
		// 从文件中读取
		config = &Configuration{
			SessionTimeout: 1,
		}
		f, err := os.Open("config.json")
		if err != nil {
			log.Fatalf("open config err: %v", err)
			return
		}
		defer f.Close()
		encoder := json.NewDecoder(f)
		err = encoder.Decode(config)
		if err != nil {
			log.Fatalf("decode config err: %v", err)
			return
		}

		// 如果环境变量有配置，读取环境变量
		ApiKey := os.Getenv("ApiKey")
		SessionTimeout := os.Getenv("SessionTimeout")
		if ApiKey != "" {
			config.ApiKey = ApiKey
		}
		if SessionTimeout != "" {
			duration, err := time.ParseDuration(SessionTimeout)
			if err != nil {
				log.Fatalf("config decode session timeout err: %v ,get is %v", err, SessionTimeout)
				return
			}
			config.SessionTimeout = duration
		}
	})
	return config
}
