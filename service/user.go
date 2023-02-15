package service

import (
	"time"

	"github.com/eryajf/chatgpt-dingtalk/config"
	"github.com/patrickmn/go-cache"
)

// UserServiceInterface 用户业务接口
type UserServiceInterface interface {
	GetUserMode(userId string) string
	SetUserMode(userId string, mode string)
	ClearUserMode(userId string)
}

var _ UserServiceInterface = (*UserService)(nil)

// UserService 用戶业务
type UserService struct {
	// 缓存
	cache *cache.Cache
}

// NewUserService 创建新的业务层
func NewUserService() UserServiceInterface {
	return &UserService{cache: cache.New(time.Second*config.LoadConfig().SessionTimeout, time.Minute*10)}
}

// GetUserMode 获取当前对话模式
func (s *UserService) GetUserMode(userId string) string {
	sessionContext, ok := s.cache.Get(userId)
	if !ok {
		return ""
	}
	return sessionContext.(string)
}

// SetUserMode 设置用户对话模式
func (s *UserService) SetUserMode(userId string, mode string) {
	s.cache.Set(userId, mode, time.Second*config.LoadConfig().SessionTimeout)
}

// ClearUserMode 重置用户对话模式
func (s *UserService) ClearUserMode(userId string) {
	s.cache.Delete(userId)
}
