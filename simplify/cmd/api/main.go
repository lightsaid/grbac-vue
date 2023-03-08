package main

import (
	"log"
	"os"
	"simplify/internal/data"
	"time"

	"github.com/gomodule/redigo/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type application struct {
	config config
	models data.Models
}

func main() {
	conf, err := loadConfig(".")
	failOnError(err, "loadConfig failed")

	db, err := openDB(conf)
	failOnError(err, "connect database failed")

	pool := openRedis(conf)

	app := application{
		config: conf,
		models: data.NewModels(db, pool),
	}

	err = app.serve()
	if err != nil {
		log.Println(err)
	}

}

func openDB(conf config) (*gorm.DB, error) {
	var colorful bool
	if conf.RunMode == "debug" {
		colorful = true
	}
	logConfig := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  colorful,    // 彩色打印
		},
	)
	return gorm.Open(mysql.Open(conf.MySQLDSN), &gorm.Config{Logger: logConfig})
}

func openRedis(conf config) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 5 * time.Minute,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", conf.RedisURL)
		},
	}
}

// failOnError 处理错误函数
func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
