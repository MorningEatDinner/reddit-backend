package mail

import (
	"fmt"
	"net/smtp"

	emailPKG "github.com/jordan-wright/email"
	"github.com/xiaorui/reddit-async/reddit-backend/settings"
	"go.uber.org/zap"
)

type SMPT struct{}

// Send: 发送邮箱验证码
func (s *SMPT) Send(email Email, config *settings.SmptConfig) bool {
	e := emailPKG.NewEmail()

	e.From = fmt.Sprintf("%v <%v>", email.From.Name, email.From.Address)
	e.To = email.To
	e.Bcc = email.Bcc
	e.Cc = email.Cc
	e.Subject = email.Subject
	e.Text = email.Text
	e.HTML = email.HTML

	err := e.Send(
		fmt.Sprintf("%v:%v", config.Host, config.Port),
		smtp.PlainAuth(
			"",
			config.Username,
			config.Password,
			config.Host,
		),
	)
	if err != nil {
		zap.L().Error("Send email failed..", zap.Error(err))
		return false
	}

	return true
}
