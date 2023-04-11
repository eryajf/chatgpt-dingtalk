package cache

import "time"

// SetAnswerID 设置用户获得答案的ID
func (s *UserService) SetAnswerID(userId, chattitle string, current uint) {
	s.cache.Set(userId+"_"+chattitle, current, time.Hour*24)
}

// GetAnswerID 获取当前用户获得答案的ID
func (s *UserService) GetAnswerID(userId, chattitle string) uint {
	sessionContext, ok := s.cache.Get(userId + "_" + chattitle)
	if !ok {
		return 0
	}
	return sessionContext.(uint)
}

// ClearUserSessionContext 清空GPT上下文，接收文本中包含 SessionClearToken
func (s *UserService) ClearAnswerID(userId, chattitle string) {
	s.cache.Delete(userId + "_" + chattitle)
}
