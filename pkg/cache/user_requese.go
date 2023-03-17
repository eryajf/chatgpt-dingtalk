package cache

import (
	"time"
)

// SetUseRequestCount 设置用户请求次数
func (s *UserService) SetUseRequestCount(userId string, current int) {
	s.cache.Set(userId+"_request", current, time.Hour*24)
}

// GetUseRequestCount 获取当前用户已请求次数
func (s *UserService) GetUseRequestCount(userId string) int {
	sessionContext, ok := s.cache.Get(userId + "_request")
	if !ok {
		return 0
	}
	return sessionContext.(int)
}
