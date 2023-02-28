package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lightsaid/grbac/errs"
	"github.com/lightsaid/grbac/helper"
	"github.com/lightsaid/grbac/initializer"
)

// 定义常量
const (
	AuthorizationKey        = "Authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

// RequireAuth 认证用户是否登录
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(AuthorizationKey)
		if len(authorizationHeader) == 0 {
			helper.ToErrResponse(c, errs.Unauthorized)
			c.Abort()
			return
		}

		// 提取 accessToken "Bearer eyJhbGciOiJIUzI...."

		// 以空格分割为两部分
		parts := strings.Fields(authorizationHeader)
		if len(parts) < 2 {
			helper.ToErrResponse(c, errs.Unauthorized.AsMessage("token 格式不匹配"))
			c.Abort()
			return
		}

		// 验证accessToken 头
		authorizationType := strings.ToLower(parts[0])
		if authorizationType != AuthorizationTypeBearer {
			helper.ToErrResponse(c, errs.Unauthorized.AsMessage("token 类型不匹配"))
			c.Abort()
			return
		}

		// 验证accessToken
		accessToken := parts[1]
		payload, err := helper.ParseToken(accessToken, initializer.App.Conf.TokenSecret)
		if err != nil {
			helper.ToErrResponse(c, errs.Unauthorized.AsMessage("token 无效"))
			c.Abort()
			return
		}

		// 设置上下文
		c.Set(AuthorizationPayloadKey, payload)

		c.Next()
	}
}
