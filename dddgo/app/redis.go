package app

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

func (app *Application) InitRedis() {
	app.RedisPool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 5 * time.Minute,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", app.Config.RedisURL)
		},
	}
}
