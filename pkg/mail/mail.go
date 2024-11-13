package mail

import (
	"sync"

	"github.com/xiaorui/reddit-async/reddit-backend/settings"
)

type From struct {
	Address string // email地址
	Name    string // 名字
}

type Email struct {
	From    From     // 发件人
	To      []string // 收件人
	Bcc     []string // 密送
	Cc      []string // 抄送
	Subject string   // 主题
	Text    []byte   // 纯文本
	HTML    []byte
}

type Mailer struct {
	Driver Driver
}

var once sync.Once
var mailer *Mailer

func NewMailer() *Mailer {
	once.Do(func() {
		mailer = &Mailer{
			Driver: &SMPT{},
		}
	})

	return mailer
}

func (m *Mailer) Send(email Email) bool {
	return m.Driver.Send(email, settings.Conf.EmailConfig.SmptConfig)
}
