package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gomodule/redigo/redis"
	"github.com/lightsaid/grbac-vue/dddgo/pkg/config"
	"github.com/lightsaid/grbac-vue/dddgo/pkg/errs"
	"github.com/lightsaid/grbac-vue/dddgo/pkg/rabbitmq"
	customerValidator "github.com/lightsaid/grbac-vue/dddgo/pkg/validator"
	"gorm.io/gorm"

	ut "github.com/go-playground/universal-translator"
)

var App *Application

type Application struct {
	Config         *config.Config
	DB             *gorm.DB
	RedisPool      *redis.Pool
	Engine         *gin.Engine
	SubPubRabbitMQ *rabbitmq.RabbitMQ
}

func NewApplication() *Application {
	conf, err := config.LoadConfig(".")
	failOnError(err, "Loading config failed")

	err = customerValidator.InitValidator()
	failOnError(err, "Init customer validator failed")

	mux := gin.Default()
	gin.SetMode(conf.RunMode)
	mq, err := rabbitmq.NewRabbitMQPubSub(conf.RegisterExchange, conf.RabbitMQURL)
	failOnError(err, "Connect rabbitmq failed")
	app := &Application{
		Config:         &conf,
		Engine:         mux,
		SubPubRabbitMQ: mq,
	}
	App = app
	return app
}

// Start 启动 App 服务 （接口服务）
func (app *Application) Start() {
	app.InitMySQL()
	app.InitRedis()
	err := app.AutoMigrate()
	failOnError(err, "AutoMigrate db failed")
}

// ToResponse 请求成功响应的处理
func ToResponse(c *gin.Context, data interface{}, msgs ...string) {
	msg := errs.Success.Message()
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	c.JSON(http.StatusOK, gin.H{
		"code": errs.Success.Code(),
		"msg":  msg,
		"data": data,
	})
}

// ToErrorResponse 请求异常的响应处理
func ToErrResponse(c *gin.Context, err *errs.AppError) {
	log.Printf("metod:%s, url: %s, error: %v\n", c.Request.Method, c.Request.URL, err)
	if err == nil {
		err = errs.ServerError
	}
	response := gin.H{
		"code": err.Code(),
		"msg":  err.Message(),
	}
	c.JSON(err.StatusCode(), response)
}

// BindRequest 如果正常绑定返回true；反之处理错误并返回false 并对请求作出响应的错误处理
func BindRequest(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBind(req); err != nil {
		return handleError(c, err)
	}

	return true
}

// BindRequestUri 绑定param参数，如：/api/users/:id 绑定id
func BindRequestUri(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindUri(req); err != nil {
		return handleError(c, err)
	}
	return true
}

// handleError 处理 Bind 请求参数发生的错误
func handleError(c *gin.Context, err error) bool {
	// 获取翻译组件实例
	v := c.Value("trans")
	trans, ok := v.(ut.Translator)
	if !ok {
		// 如果trans经过setTranslations中间件设置后还获取不到，那就是服务内部出问题
		ToErrResponse(c, errs.ServerError.AsException(err))
		return false
	}
	// 断言错误是否为 validator/v10 的验证错误信息
	verrs, ok := err.(validator.ValidationErrors)
	if !ok { // 其他方面的参数不匹配
		ToErrResponse(c, errs.BadRequest.AsException(err))
		return false
	}

	// 对错误信息进行翻译, 得到的是 map[string]string 结构数据
	merrs := verrs.Translate(trans)

	// 拼接错误消息
	var msg string
	var index = 0
	for _, v := range merrs {
		if index > 0 {
			msg += ";"
		}
		msg += v
		index++
	}

	ToErrResponse(c, errs.BadRequest.AsException(err, msg))

	return false
}

// failOnError 处理错误函数
func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
