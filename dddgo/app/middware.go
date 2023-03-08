package app

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/lightsaid/grbac-vue/dddgo/pkg/errs"
	"github.com/lightsaid/grbac-vue/dddgo/pkg/token"
)

// 定义常量
const (
	AuthorizationKey        = "Authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

// Translation 设置翻译 ut.Translator 实例上下文
func (app *Application) Translation() gin.HandlerFunc {
	return func(c *gin.Context) {
		uni := ut.New(en.New(), zh.New(), zh_Hant_TW.New())
		locale := c.GetHeader("locale")
		trans, _ := uni.GetTranslator(locale)
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if ok {
			switch locale {
			case "zh":
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
			case "en":
				_ = en_translations.RegisterDefaultTranslations(v, trans)
			default:
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
			}
			c.Set("trans", trans)
		}
		c.Next()
	}
}

// RequireAuth 认证用户是否登录
func (app *Application) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(AuthorizationKey)
		if len(authorizationHeader) == 0 {
			ToErrResponse(c, errs.Unauthorized)
			c.Abort()
			return
		}

		// 提取 accessToken "Bearer eyJhbGciOiJIUzI...."

		// 以空格分割为两部分
		parts := strings.Fields(authorizationHeader)
		if len(parts) < 2 {
			ToErrResponse(c, errs.Unauthorized.AsMessage("token 格式不匹配"))
			c.Abort()
			return
		}

		// 验证accessToken 头
		authorizationType := strings.ToLower(parts[0])
		if authorizationType != AuthorizationTypeBearer {
			ToErrResponse(c, errs.Unauthorized.AsMessage("token 类型不匹配"))
			c.Abort()
			return
		}

		// 验证accessToken
		accessToken := parts[1]
		payload, err := token.ParseToken(accessToken, app.Config.TokenSecret)
		if err != nil {
			ToErrResponse(c, errs.Unauthorized.AsMessage("token 无效"))
			c.Abort()
			return
		}

		// 设置上下文
		c.Set(AuthorizationPayloadKey, payload)

		c.Next()
	}
}
