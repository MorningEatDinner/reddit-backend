package redis

import "time"

// SetVerifyCode： 将手机号码的验证码存储到redis中
func SetVerifyCode(phone, code string) error {
	expireTime := time.Minute * 5

	return RDB.Client.Set(RDB.Context, getRedisKey(KeyVerifyCode)+phone, code, expireTime).Err()
}

// CheckVerifyCode： 验证验证是否正确
func CheckVerifyCode(key, answer string) bool {
	val := RDB.Client.Get(RDB.Context, getRedisKey(KeyVerifyCode)+key).Val()
	return val == answer
}
