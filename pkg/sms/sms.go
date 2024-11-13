package sms

import (
	"sync"

	"github.com/xiaorui/reddit-async/reddit-backend/settings"
)

type SMS struct {
	Driver Driver
}

var once sync.Once
var sms *SMS

func NewSms() *SMS {
	once.Do(func() {
		sms = &SMS{Driver: &Aliyun{}}
	})

	return sms
}

func (sms *SMS) Send(phone, message string) bool {
	return sms.Driver.Send(phone, message, settings.Conf.SmsConfig)
}
