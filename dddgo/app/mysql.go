package app

import (
	"log"
	"os"
	"time"

	"github.com/lightsaid/grbac-vue/dddgo/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitMySQL 初始化数据库链接
func (app *Application) InitMySQL() {
	var err error
	logConfig := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 彩色打印
		},
	)

	app.DB, err = gorm.Open(mysql.Open(app.Config.MySQLDSN), &gorm.Config{Logger: logConfig})
	if err != nil {
		log.Fatal("连接 MySQL 数据库失败 ", err.Error())
	}
}

// AutoMigrate 自动迁移表
func (app *Application) AutoMigrate() error {
	err := app.DB.AutoMigrate(&entity.User{}, &entity.Role{}, &entity.Permission{})
	return err
}
