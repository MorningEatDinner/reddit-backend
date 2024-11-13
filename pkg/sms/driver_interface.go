package sms

import "github.com/xiaorui/reddit-async/reddit-backend/settings"

// 发送短信的接口 脱离具体形态
type Driver interface {
	// 发送短信
	Send(phone, message string, config *settings.SmsConfig) bool
}
