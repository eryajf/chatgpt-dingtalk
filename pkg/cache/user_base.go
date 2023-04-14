package cache

import (
	"time"

	"github.com/eryajf/chatgpt-dingtalk/public"
	"github.com/patrickmn/go-cache"
)

// UserServiceInterface 用户业务接口
type UserServiceInterface interface {
	// 用户聊天模式
	GetUserMode(userId string) string
	SetUserMode(userId, mode string)
	ClearUserMode(userId string)
	// 用户聊天上下文
	GetUserSessionContext(userId string) string
	SetUserSessionContext(userId, content string)
	ClearUserSessionContext(userId string)
	// 用户请求次数
	SetUseRequestCount(userId string, current int)
	GetUseRequestCount(uerId string) int
	// 用户对话ID
	SetAnswerID(userId, chattype string, current uint)
	GetAnswerID(uerId, chattype string) uint
	ClearAnswerID(userId, chattitle string)
}

var _ UserServiceInterface = (*UserService)(nil)

// UserService 用戶业务
type UserService struct {
	// 缓存
	cache *cache.Cache
}

// NewUserService 创建新的业务层
func NewUserService() UserServiceInterface {
	return &UserService{cache: cache.New(time.Second*public.Config.SessionTimeout, time.Hour*1)}
}
