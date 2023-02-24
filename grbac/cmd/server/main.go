package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lightsaid/grbac/docs"
	"github.com/lightsaid/grbac/initializer"
	"github.com/lightsaid/grbac/routers"
)

// @title 权限管理系统接口文档
// @version 1.0
// @description 这是一个RBAC权限管理系统接口文档，访问地址是: localhost:9999/docs/swagger/index.html
// @termsOfService https://github.com/lightsaid/grbac-vue

// @host localhost:9999
// @BasePath /v1/api
func main() {
	var err error
	initializer.NewAppConfig(".")
	initializer.InitMySQL()
	err = initializer.AutoMigrate()
	failOnError(err, "initializer.AutoMigrate()")
	err = initializer.InitValidator()
	failOnError(err, "initializer.InitValidator()")
	gin.SetMode(initializer.App.Conf.RunMode)
	engine := gin.Default()
	routers.SetupRoutes(engine)

	server := http.Server{
		Addr:           initializer.App.Conf.HTTPServerAddress,
		Handler:        engine,
		IdleTimeout:    time.Minute,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 4 << 20, // 4M
	}

	// 启动服务监听
	go func() {
		log.Println("Starting server on ", initializer.App.Conf.HTTPServerAddress)
		if err := server.ListenAndServe(); err != nil {
			log.Println("ListenAndServe: ", err)
		}
	}()

	// 优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	//  阻塞，等待 os.Signal 信号
	<-quit
	log.Println("Stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown server error: ", err)
	}
	log.Println("Stopped server.")
}

// failOnError 打印错误，终止程序
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
