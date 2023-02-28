package initializer

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

var RedisPool *redis.Pool

func InitRedis() {
	RedisPool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 5 * time.Minute,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", App.Conf.RedisURL)
		},
	}
}
