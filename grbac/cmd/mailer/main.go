package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/lightsaid/grbac/initializer"
	"github.com/lightsaid/grbac/mailer"
	"github.com/lightsaid/grbac/models"
	"github.com/lightsaid/grbac/rabbitmq"
)

func main() {
	var err error
	initializer.NewAppConfig(".")
	mq, err := rabbitmq.NewRabbitMQPubSub(
		initializer.App.Conf.RegisterExchange,
		initializer.App.Conf.RabbitMQURL,
	)

	if err != nil {
		log.Fatal(err)
	}

	var callback = func(msg string) {
		// TODO:
		fmt.Println("Mailer Server 接收到 Message: ", msg)
		// 发送邮件
		sender := mailer.NewGmailSender(
			initializer.App.Conf.MailSenderName,
			initializer.App.Conf.MailSenderAddress,
			initializer.App.Conf.MailSenderPassword,
		)

		var payload models.RegisterMailerPaylod
		err := json.Unmarshal([]byte(msg), &payload)
		if err != nil {
			log.Println("解析 RegisterMailerPaylod 错误： ", err)
			return
		}

		err = sender.SendEmail(
			"账户激活",
			payload.Content,
			[]string{payload.Email},
			nil,
			nil,
			nil,
		)
		if err != nil {
			log.Println("发送邮件错误： ", err)
			return
		}
	}

	var errChan = make(chan error)
	defer close(errChan)

	go func() {
		// 等待是否有错误
		err = <-errChan
		if err != nil {
			panic(err)
		}
	}()

	// 执行消费
	mq.ConsumerPubSubCtx(errChan, callback)
}
