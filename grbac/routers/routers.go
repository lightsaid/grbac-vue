package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lightsaid/grbac/controllers"
)

func InitRoutes(engine *gin.Engine) {
	// 健康检查路由
	engine.GET("/v1/api/health", controllers.CheckHealth)

	// 发送邮件
	engine.POST("/v1/api/sendemail", controllers.SendGoMail)

	// auth 认证相关路由
	auth := engine.Group("/v1/api/auth")
	auth.POST("/register", controllers.Register)                   // 注册
	auth.GET("/verifyemail/:verifyCode", controllers.ActivateUser) // 激活用户
	auth.POST("/login", controllers.Login)                         // 登入
	auth.POST("/logout", controllers.Logout)                       // 登出
	auth.POST("/refresh", controllers.Refresh)                     // 刷新 Token
	auth.POST("/forgotpswd", controllers.ForgotPswd)               // 忘记密码
	auth.PATCH("/restpswd/:restToken", controllers.RestPswd)       // 重置密码

	// admin 管理员路由
	admin := engine.Group("/v1/api/admin")
	admin.GET("/users/profile", controllers.GetProfile)        // 获取自己个人信息
	admin.GET("/users/list", controllers.ListUsers)            // 获取用户列表 /users/list?page=1&size=10
	admin.PUT("/users/update/:uid", controllers.UpdateProfile) // 更新用户信息

}
