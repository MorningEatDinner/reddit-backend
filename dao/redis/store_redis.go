package redis

import (
	"errors"
	"time"
)

/*
	实现captcha的redis store
*/

type RedisStore struct {
	RedisClient *RedisClient
}

// Set： 设置验证码id和答案
func (s *RedisStore) Set(id string, value string) error {
	expireTime := time.Minute * 15 // 后面可以在配置中导入
	if err := s.RedisClient.Client.Set(s.RedisClient.Context, getRedisKey(KeyCaptcha)+id, value, expireTime).Err(); err != nil {
		return errors.New("存储数据失败")
	}
	return nil
}

// Get：获得验证码的值
func (s *RedisStore) Get(id string, clear bool) string {
	val := s.RedisClient.Client.Get(s.RedisClient.Context, getRedisKey(KeyCaptcha)+id).Val()
	if clear {
		s.RedisClient.Client.Del(s.RedisClient.Context, id)
	}
	return val
}

// Verify：验证验证码是否正确
func (s *RedisStore) Verify(id, answer string, clear bool) bool {
	realVal := s.Get(id, clear)
	return realVal == answer
}
