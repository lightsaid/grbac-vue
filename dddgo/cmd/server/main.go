package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lightsaid/grbac-vue/dddgo/app"
	v1 "github.com/lightsaid/grbac-vue/dddgo/app/v1"
	_ "github.com/lightsaid/grbac-vue/dddgo/docs"
)

// @title 权限管理系统接口文档
// @version 1.0
// @description 这是一个RBAC权限管理系统接口文档，访问地址是: localhost:9999/docs/swagger/index.html
// @termsOfService https://github.com/lightsaid/grbac-vue

// @host localhost:9999
// @BasePath /v1/api
func main() {
	app := app.NewApplication()
	app.Start()
	v1.InitRoutes()
	server := http.Server{
		Addr:           app.Config.HTTPServerAddress,
		Handler:        app.Engine,
		IdleTimeout:    time.Minute,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 4 << 20, // 4M
	}

	// 启动服务监听
	go func() {
		log.Println("Starting server on ", app.Config.HTTPServerAddress)
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

	// 释放资源
	app.SubPubRabbitMQ.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 关机
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown server error: ", err)
	}
	log.Println("Stopped server.")
}
