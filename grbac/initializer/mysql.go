package initializer

import (
	"log"
	"os"
	"time"

	"github.com/lightsaid/grbac/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 提供 DB 全局使用
var DB *gorm.DB

// InitMySQL 初始化数据库链接
func InitMySQL() {
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

	DB, err = gorm.Open(mysql.Open(os.Getenv("DSN")), &gorm.Config{Logger: logConfig})
	if err != nil {
		log.Fatal("连接 MySQL 数据库失败! ", err.Error())
	}
}

// AutoMigrate 自动迁移表
func AutoMigrate() error {
	err := DB.AutoMigrate(models.User{})
	return err
}
