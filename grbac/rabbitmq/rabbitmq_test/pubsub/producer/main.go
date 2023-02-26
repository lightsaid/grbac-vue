package main

import (
	"context"
	"fmt"

	"github.com/lightsaid/grbac/rabbitmq"
)

const mqsource = "amqp://guest:guest@localhost:5666/"

func main() {
	mq, err := rabbitmq.NewRabbitMQPubSub("logs", mqsource)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 10; i++ {
		mq.PublishPubSubCtx(context.Background(), fmt.Sprintf("第 %d 条消息", i+1))
	}
}
