package cache

import "github.com/patrickmn/go-cache"

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
