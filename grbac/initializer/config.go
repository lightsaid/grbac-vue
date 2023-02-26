package initializer

import (
	"log"
	"sync"
	"time"

	"github.com/lightsaid/grbac/rabbitmq"
	"github.com/spf13/viper"
)

var App *appConfig

type config struct {
	RunMode              string        `mapstructure:"RUN_MODE"`             // HTTP Server 启动模式: debug | release
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`  // HTTP Server IP地址+端口
	MySQLDSN             string        `mapstructure:"MYSQL_DSN"`            // 链接 MySQL 的 DSN
	MailSenderName       string        `mapstructure:"MAIL_SENDER_NAME"`     // 发送邮件人的名字，对方收到邮件会显示
	MailSenderAddress    string        `mapstructure:"MAIL_SENDER_ADDRESS"`  // 发送人邮件人的邮箱地址
	MailSenderPassword   string        `mapstructure:"MAIL_SENDER_PASSWORD"` // 发送人邮件的应用专用密码
	To163MailAddress     string        `mapstructure:"TO_163_MAIL_ADDRESS"`  // 测试163邮箱接收邮件的邮箱地址
	RabbitMQURL          string        `mapstructure:"RABBITMQ_URL"`
	ActivateEmailURL     string        `mapstructure:"ACTIVATE_EMAIL_URL"`
	SignatureSecret      string        `mapstructure:"SIGNATURE_SECRET"`
	RegisterExchange     string        `mapstructure:"REGISTER_EXCHANGE"`
	TokenSecret          string        `mapstructure:"TOKEN_SECRET"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESSTOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESHTOKEN_DURATION"`
}

// appConfig 定义一个结构体保存全局配置
type appConfig struct {
	Conf *config

	SubPubRabbitMQ *rabbitmq.RabbitMQ

	// 其他配置
	Wait *sync.WaitGroup
}

// NewAppConfig 加载 app.env 配置文件到 appConfig 结构体，定义其他配置
// 得到一个全局的配置变量 App
func NewAppConfig(path string) {
	conf, err := loadConfig(path)
	if err != nil {
		log.Fatalf("Error loading app.env file %s", err)
	}

	mq, err := rabbitmq.NewRabbitMQPubSub(conf.RegisterExchange, conf.RabbitMQURL)
	if err != nil {
		log.Fatal("rabbitmq.NewRabbitMQPubSub: ", err)
	}

	wg := sync.WaitGroup{}

	App = &appConfig{
		Conf:           &conf,
		Wait:           &wg,
		SubPubRabbitMQ: mq,
	}
}

// LoadConfig 从配置文件和环境读取配置参数
func loadConfig(path string) (conf config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&conf)
	return
}
