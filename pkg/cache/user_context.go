package cache

import "github.com/patrickmn/go-cache"

// SetUserSessionContext 设置用户会话上下文文本，question用户提问内容，GPT回复内容
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

// ClearUserSessionContext 清空GPT上下文，接收文本中包含 SessionClearToken
func (s *UserService) ClearUserSessionContext(userId string) {
	s.cache.Delete(userId + "_content")
}
