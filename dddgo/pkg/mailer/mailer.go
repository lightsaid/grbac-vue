package mailer

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

// 此包实现发送Gmail邮件功能；而QQ邮件或163邮件类似

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type Mailer interface {
	// SendEmail 发送邮件接口
	// subject 主题、content 内容、to 发给谁、cc 抄送给谁、bcc 密件抄送、attachFiles 附件
	SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error
}

type GmailSender struct {
	name              string // 发件人
	fromEmailAddress  string // 发件人邮箱地址
	fromEmailPassword string // 发件人邮箱授权密码
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) Mailer {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

// SendEmail 发送邮件接口
// subject 主题、content 内容、to 发给谁、cc 抄送给谁、bcc 密件抄送、attachFiles 附件
func (sender *GmailSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
	mail := email.NewEmail()
	mail.From = sender.name + "<" + sender.fromEmailAddress + ">"
	mail.Subject = subject
	mail.HTML = []byte(content)
	mail.To = to
	mail.Cc = cc
	mail.Bcc = bcc

	for _, f := range attachFiles {
		_, err := mail.AttachFile(f)
		if err != nil {
			return fmt.Errorf("加载附近出错,请确保路径正确 %s: %w", f, err)
		}
	}

	// identity 通常是空字符串
	auth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return mail.Send(smtpServerAddress, auth)
}
