package helper

import (
	"log"
	"net/http"

	"github.com/lightsaid/grbac/errs"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// ToResponse 请求成功响应的处理
func ToResponse(c *gin.Context, data interface{}) {
	msg, ok := data.(string)
	if ok {
		c.JSON(http.StatusOK, gin.H{
			"code": errs.StatusOK.Code(),
			"msg":  msg,
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": errs.StatusOK.Code(),
		"msg":  errs.StatusOK.Message(),
		"data": data,
	})
}

// ToErrorResponse 请求异常的响应处理
func ToErrResponse(c *gin.Context, err *errs.AppError) {
	log.Printf("metod:%s, url: %s, error: %v\n", c.Request.Method, c.Request.URL, err)
	if err == nil {
		err = errs.InternalServerError
	}
	response := gin.H{
		"code": err.Code(),
		"msg":  err.Message(),
	}
	c.JSON(err.StatusCode(), response)
}

// BindRequest 如果正常绑定返回true；反之处理错误并返回false 并对请求作出响应的错误处理
func BindRequest(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
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
		ToErrResponse(c, errs.InternalServerError.AsException(err))
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
