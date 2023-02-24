package mailer

import (
	"os"
	"testing"

	"github.com/lightsaid/grbac/initializer"
	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	initializer.InitConfig("../.env")
	sender := NewGmailSender(os.Getenv("EMAIL_SENDER_NAME"), os.Getenv("EMAIL_SENDER_ADDRESS"), os.Getenv("EMAIL_SENDER_PASSWORD"))

	subject := "Test Send Email"
	content := `
		<h1>Hello World! 欢迎～
		<p>欢迎查看我的a href="https://github.com/lightsaid/grbac-vue">Github<</a></p>
	`
	// 测试发送到gmail和163
	to := []string{os.Getenv("EMAIL_SENDER_ADDRESS"), os.Getenv("EMAIL_163_SENDER_ADDRESS")}
	attchFiles := []string{"../../README.md"}

	err := sender.SendEmail(subject, content, to, nil, nil, attchFiles)
	require.NoError(t, err)

}
