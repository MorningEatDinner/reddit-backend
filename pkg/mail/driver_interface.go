package mail

import "github.com/xiaorui/reddit-async/reddit-backend/settings"

type Driver interface {
	// 发送验证码
	Send(email Email, config *settings.SmptConfig) bool
}
