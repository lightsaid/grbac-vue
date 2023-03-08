package v1

import (
	"github.com/lightsaid/grbac-vue/dddgo/app"
	"github.com/lightsaid/grbac-vue/dddgo/domain"
	"github.com/lightsaid/grbac-vue/dddgo/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoutes() {
	mux := app.App.Engine

	mux.Use(app.App.Translation())

	otherHandler := OtherHandler{}
	mux.GET("/v1/api/health", otherHandler.CheckHealthHandler)

	// 注册Swagger文档路由
	mux.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	authHandler := AuthHandler{
		service: service.NewAuthService(
			domain.NewUserRepository(app.App.DB),
			domain.NewAuthRepository(app.App.DB),
			domain.NewSessionRepository(app.App.RedisPool)),
	}

	// auth 认证相关路由
	auth := mux.Group("/v1/api/auth")
	auth.POST("/register", authHandler.RegisterHandler)                   // 注册
	auth.GET("/verifyemail/:verifyCode", authHandler.ActivateUserHandler) // 激活用户
	auth.POST("/login", authHandler.LoginHandler)                         // 登入
	auth.POST("/refresh", authHandler.RefreshHandler)                     // 刷新 Token

	admin := mux.Group("/v1/api/admin").Use(app.App.RequireAuth())
	admin.GET("/users/:id") // 获取一个用户信息
	admin.GET("/users")     // 获取用户列表
	admin.PUT("/users/:id") // 更新一个用户信息

	admin.POST("/roles")       // 添加一个角色
	admin.PUT("/roles/:id")    // 更新一个角色
	admin.GET("/roles")        // 获取角色列表
	admin.GET("/roles/:id")    // 获取一个角色
	admin.DELETE("/roles/:id") // 删除一个角色

	admin.POST("/premissions")    // 添加一个权限
	admin.PUT("/premissions/:id") // 更新一个权限
	admin.GET("/premissions")     // 获取权限列表
	admin.GET("/premissions/:id") // 获取一个权限
	admin.GET("/premissions/:id") // 删除一个权限

}
