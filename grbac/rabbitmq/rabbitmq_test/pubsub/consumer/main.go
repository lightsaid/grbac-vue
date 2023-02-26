package main

import (
	"fmt"
	"log"

	"github.com/lightsaid/grbac/rabbitmq"
)

const mqsource = "amqp://guest:guest@localhost:5666/"

func main() {
	mq, err := rabbitmq.NewRabbitMQPubSub("logs", mqsource)
	if err != nil {
		panic(err)
	}

	var callback = func(msg string) {
		log.Println("CallBack Get Message: ", msg)
	}

	var errChan = make(chan error)
	defer close(errChan)
	defer func() {
		fmt.Println("关闭 errChann")
	}()

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
