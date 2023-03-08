package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/lightsaid/grbac-vue/dddgo/app"
	"github.com/lightsaid/grbac-vue/dddgo/dto"
	"github.com/lightsaid/grbac-vue/dddgo/pkg/mailer"
	"github.com/lightsaid/grbac-vue/dddgo/pkg/rabbitmq"
)

func main() {
	var err error
	app := app.NewApplication()
	mq, err := rabbitmq.NewRabbitMQPubSub(
		app.Config.RegisterExchange,
		app.Config.RabbitMQURL,
	)

	if err != nil {
		log.Fatal(err)
	}

	var callback = func(msg string) {
		// TODO:
		fmt.Println("Mailer Server 接收到 Message: ", msg)
		// 发送邮件
		sender := mailer.NewGmailSender(
			app.Config.MailSenderName,
			app.Config.MailSenderAddress,
			app.Config.MailSenderPassword,
		)

		var payload dto.RegisterMailerPaylod
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
