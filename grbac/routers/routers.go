package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lightsaid/grbac/controllers"
	"github.com/lightsaid/grbac/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(engine *gin.Engine) {
	// 中间件设置
	engine.Use(middleware.Translation())

	// 健康检查路由
	engine.GET("/v1/api/health", controllers.CheckHealth)

	// 发送邮件，转由 mailer server 实现，通过 RabbitMQ 订阅发布
	// engine.POST("/v1/api/sendemail", controllers.SendGoMail)

	// 注册Swagger文档路由
	engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// auth 认证相关路由
	auth := engine.Group("/v1/api/auth")
	auth.POST("/register", controllers.Register)                   // 注册
	auth.GET("/verifyemail/:verifyCode", controllers.ActivateUser) // 激活用户
	auth.POST("/login", controllers.Login)                         // 登入
	auth.POST("/refresh", controllers.Refresh)                     // 刷新 Token

	// 上传文件
	engine.POST("/v1/api/upload", controllers.UploadFiles).Use(middleware.RequireAuth())

	// admin 管理员路由
	admin := engine.Group("/v1/api/admin").Use(middleware.RequireAuth())
	{
		// User 模块
		admin.GET("/users/profile", controllers.GetProfile)        // 获取自己个人信息
		admin.GET("/users/list", controllers.ListUsers)            // 获取用户列表 /users/list?page=1&size=10
		admin.PUT("/users/update/:uid", controllers.UpdateProfile) // 更新用户信息
		admin.POST("/users/logout", controllers.Logout)            // 登出

		// Role 角色模块
		admin.POST("/roles", controllers.CreateRole)       // 创建角色
		admin.GET("/roles", controllers.ListRoles)         // 获取角色列表 roles?page=1&size=10
		admin.PUT("/roles/:id", controllers.UpdateRole)    // 更新角色
		admin.GET("/roles/:id", controllers.GetRole)       // 获取一个角色
		admin.DELETE("/roles/:id", controllers.DeleteRole) // 删除角色

		// Permission 权限模块
		admin.POST("/permissions", controllers.CreatePermission)       // 创建一个权限
		admin.GET("/permissions", controllers.ListPermissions)         // 获取权限列表 /permissions?page=1&size=10
		admin.PUT("/permissions/:id", controllers.UpdatePermission)    // 更新一个权限
		admin.GET("/permissions/:id", controllers.GetPermission)       // 获取一个权限
		admin.DELETE("/permissions/:id", controllers.DeletePermission) // 删除一个权限

	}

}
