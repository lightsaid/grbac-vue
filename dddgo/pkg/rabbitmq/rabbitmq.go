package rabbitmq

import (
	"context"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	// 链接
	conn *amqp091.Connection
	// channel
	ch *amqp091.Channel
	// 队列名称
	QueueName string
	// 交换机
	Exchange string
	// key
	Key string
	// 链接信息 连接到rabbitmq服务的地址
	MQSource string
}

// failOnError 处理错误函数
func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
		// panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// NewRabbitMQ 创建一个RabbitMQ结构体实例, 包括连接到rabbitmq服务，获取channel
func NewRabbitMQ(queueName string, exchange string, key string, mqSource string) (*RabbitMQ, error) {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		MQSource:  mqSource,
	}

	var err error

	// 1. 连接到rabbitmq服务
	rabbitmq.conn, err = amqp091.Dial(rabbitmq.MQSource)
	failOnError(err, "创建连接错误")
	if err != nil {
		return nil, err
	}

	// 2. 获取一个channel
	rabbitmq.ch, err = rabbitmq.conn.Channel()
	failOnError(err, "获取Channel失败")
	if err != nil {
		return nil, err
	}
	return rabbitmq, nil
}

// Close 释放资源，断开channel和connection
func (r *RabbitMQ) Close() {
	r.ch.Close()
	r.conn.Close()
}

// ================ 发布订阅模式 Start================

// NewRabbitMQPubSub 创建订阅模式实例方法
func NewRabbitMQPubSub(exchange string, source string) (*RabbitMQ, error) {
	return NewRabbitMQ("", exchange, "", source)
}

// PublishPubSubCtx 发布订阅模式发送消息
func (r *RabbitMQ) PublishPubSubCtx(ctx context.Context, msg string) error {
	// 1. 声明交换机
	err := r.ch.ExchangeDeclare(
		r.Exchange, // name
		"fanout",   // type 类型， 订阅模式必须是fanout
		true,       // durable 持久化
		false,      // auto-deleted 是否自动删除
		false,      // 内部使用（当前程序使用，其他端使用不了）
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "r.ch.ExchangeDeclare error, exchangeName: "+r.Exchange)

	// 2. 发送消息
	err = r.ch.PublishWithContext(
		ctx,
		r.Exchange,
		"",
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
	failOnError(err, "r.ch.PublishWithContext error, exchangeName: "+r.Exchange)
	if err != nil {
		return err
	}
	return nil
}

// ConsumerPubSubCtx 订阅发布模式下发送消息
func (r *RabbitMQ) ConsumerPubSubCtx(errChan chan error, fn func(msg string)) {
	// 1. 声明交换机
	err := r.ch.ExchangeDeclare(
		r.Exchange, // name
		"fanout",   // type 类型， 订阅模式必须是fanout
		true,       // durable 持久化
		false,      // auto-deleted 是否自动删除
		false,      // 内部使用（当前程序使用，其他端使用不了）
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "r.ch.ExchangeDeclare error, exchangeName: "+r.Exchange)

	// 2. 声明队列
	q, err := r.ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "r.ch.QueueDeclare")

	// 3. 队列绑定交换机
	err = r.ch.QueueBind(q.Name, "", r.Exchange, false, nil)
	failOnError(err, "r.ch.QueueBind")
	if err != nil {
		errChan <- err
	}

	// 4. 消费消息
	msgs, err := r.ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "r.ch.Consume")
	if err != nil {
		errChan <- err
	}

	// 如果能执行到这里已经没有错误了
	errChan <- nil

	forever := make(chan struct{})
	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s\n", string(msg.Body))
			fn(string(msg.Body))
		}
	}()
	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	// 阻塞
	<-forever
}

// ================ 发布订阅模式模式 End================

// ================ 简单模式 Start================

// NewRabbitMQSimple 创建Simple模式实例的方法，仅需要queueName和链接到rabbitmq服务地址
func NewRabbitMQSimple(queueName string, source string) (*RabbitMQ, error) {
	return NewRabbitMQ(queueName, "", "", source)
}

// PublishSimpleCtx Simple 生产者模式下发送消息
func (r *RabbitMQ) PublishSimpleCtx(ctx context.Context, msg string) error {
	// 1. 声明队列，队列不存在创建，存在则获取
	_, err := r.ch.QueueDeclare(
		r.QueueName, // name 队列名
		false,       // durable 是否持久化
		false,       // autoDelete 是否自动删除，当最后一个消费者连接断开，是否删除
		false,       // exclusive 是否具有排他性（独占），也就是只能被当前客户端使用，其他消费者使用不了
		false,       // noWait 是否阻塞，发送消息后是否等待服务器有响应
		nil,         // 额外参数
	)
	if err != nil {
		return err
	}

	// 2. 发送消息
	err = r.ch.PublishWithContext(
		ctx,
		r.Exchange,  // 默认 direct
		r.QueueName, // key 默认使用 QueueName
		false,       // mandatory 如果为True，根据exchange类型和routing key规则，如果无法找到符合规定的队列，就会把发送的消息会退
		false,       // immediate 如果为True，当exchange发送消息到队列后发现队列没绑定消费者，则把消息会退给发送者
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
	return err
}

// ConsumeSimple 简单模式消费消息
func (r *RabbitMQ) ConsumeSimple() {
	// 1. 声明队列，队列不存在创建，存在则获取
	_, err := r.ch.QueueDeclare(
		r.QueueName, // name 队列名
		false,       // durable 是否持久化
		false,       // autoDelete 是否自动删除，当最后一个消费者连接断开，是否删除
		false,       // exclusive 是否具有排他性（独占），也就是只能被当前客户端使用，其他消费者使用不了
		false,       // noWait 是否阻塞，发送消息后是否等待服务器有响应
		nil,         // 额外参数
	)
	if err != nil {
		log.Println(err)
	}

	// 2. 接收消息
	msgs, err := r.ch.Consume(
		r.QueueName,
		"",    // consumer 用来区分消费者
		true,  // autoAck 是否自动应答
		false, // exclusive 是否具有排他性(独占)
		false, // noLocal 如果设置为True，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false, // noWait 队列消费是否阻塞，false 表示阻塞，消费完这个消息再来下一个，一个一个来
		nil,   // args
	)
	if err != nil {
		log.Println(err)
	}

	forever := make(chan struct{})

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s\n", string(msg.Body))
		}
	}()
	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	// 阻塞
	<-forever
}

// ================ 简单模式 End================
