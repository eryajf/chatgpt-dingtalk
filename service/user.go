package service

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// UserServiceInterface 用户业务接口
type UserServiceInterface interface {
	GetUserMode(userId string) string
	SetUserMode(userId, mode string)
	ClearUserMode(userId string)
	GetUserSessionContext(userId string) string
	SetUserSessionContext(userId, content string)
	ClearUserSessionContext(userId string)
}

var _ UserServiceInterface = (*UserService)(nil)

// UserService 用戶业务
type UserService struct {
	// 缓存
	cache *cache.Cache
}

// NewUserService 创建新的业务层
func NewUserService() UserServiceInterface {
	return &UserService{cache: cache.New(time.Hour*2, time.Hour*5)}
}

// GetUserMode 获取当前对话模式
func (s *UserService) GetUserMode(userId string) string {
	sessionContext, ok := s.cache.Get(userId + "_mode")
	if !ok {
		return ""
	}
	return sessionContext.(string)
}

// SetUserMode 设置用户对话模式
func (s *UserService) SetUserMode(userId string, mode string) {
	s.cache.Set(userId+"_mode", mode, cache.DefaultExpiration)
}

// ClearUserMode 重置用户对话模式
func (s *UserService) ClearUserMode(userId string) {
	s.cache.Delete(userId + "_mode")
}

// SetUserSessionContext 设置用户会话上下文文本，question用户提问内容，GTP回复内容
func (s *UserService) SetUserSessionContext(userId string, content string) {
	s.cache.Set(userId+"_content", content, cache.DefaultExpiration)
}

// GetUserSessionContext 获取用户会话上下文文本
func (s *UserService) GetUserSessionContext(userId string) string {
	sessionContext, ok := s.cache.Get(userId + "_content")
	if !ok {
		return ""
	}
	return sessionContext.(string)
}

// ClearUserSessionContext 清空GTP上下文，接收文本中包含 SessionClearToken
func (s *UserService) ClearUserSessionContext(userId string) {
	s.cache.Delete(userId + "_content")
}
