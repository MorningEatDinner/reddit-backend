package redis

//存放redis key

const (
	KeyPrefix          = "bluebell:"
	KeyPostTimeZSet    = "post:time:"
	KeyPostScoreZSet   = "post:score:"
	KeyPostVotedZSetPF = "post:voted:"            // 这里更改了好像会有点麻烦
	KeyCommunitySetPF  = "community:"             // 保存每个community下面的post的集合
	KeyCaptcha         = "signup:captcha:"        // 保存图形验证码
	KeyVerifyCode      = "signup:verifycode:"     // 保存短信或邮件验证码
	KeyProfileStatus   = "signup:profile_status:" // 是否更新了个人信息
)

// 给key加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}

func GetRedisKey(key string) string {
	return getRedisKey(key)
}
