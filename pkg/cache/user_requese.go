package cache

import (
	"time"
)

// SetUseRequestCount 设置用户请求次数
func (s *UserService) SetUseRequestCount(userId string, current int) {
	expiration := time.Now().Add(time.Hour * 24).Truncate(time.Hour * 24)
	duration := expiration.Sub(time.Now())
	// 设置缓存失效时间为第二天零点
	s.cache.Set(userId+"_request", current, duration)
}

// GetUseRequestCount 获取当前用户已请求次数
func (s *UserService) GetUseRequestCount(userId string) int {
	sessionContext, ok := s.cache.Get(userId + "_request")
	if !ok {
		return 0
	}
	return sessionContext.(int)
}
