package mailer

import (
	"testing"

	"github.com/lightsaid/grbac/initializer"
	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	initializer.NewAppConfig("../")
	sender := NewGmailSender(
		initializer.App.Conf.MailSenderName,
		initializer.App.Conf.MailSenderAddress,
		initializer.App.Conf.MailSenderPassword,
	)

	subject := "Test Send Email"
	content := `
		<h1>Hello World! 欢迎～
		<p>欢迎查看我的a href="https://github.com/lightsaid/grbac-vue">Github<</a></p>
	`
	// 测试发送到gmail和163
	to := []string{initializer.App.Conf.MailSenderAddress, initializer.App.Conf.To163MailAddress}
	attchFiles := []string{"../../README.md"}

	err := sender.SendEmail(subject, content, to, nil, nil, attchFiles)
	require.NoError(t, err)

}
