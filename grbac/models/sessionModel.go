package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var ErrSessionNotFound = errors.New("not found")

const (
	timeLayout = "2006-01-02 15:04:05"
)

// Session 登录 Session
type Session struct {
	// token id
	TID string `redis:"tid" json:"tid"`
	// 用户 id
	UID          uint   `redis:"uid" json:"uid"`
	RefreshToken string `redis:"refresh_token" json:"refresh_token"`
	// 设置 redis 不扫描，以字符串格式存储
	ExpiresAt time.Time `redis:"_" json:"expires_at"`
	CreatedAt time.Time `redis:"_" json:"created_at"`
	UserAgent string    `redis:"user_agent" json:"user_agent"`
	ClientIP  string    `redis:"client_ip" json:"client_ip"`
}

// Save 保存
func (s *Session) Save(pool *redis.Pool, key string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do(
		"HMSET", key,
		"tid", s.TID,
		"uid", s.UID,
		"refresh_token", s.RefreshToken,
		"expires_at", s.ExpiresAt.Format(timeLayout),
		"created_at", s.CreatedAt.Format(timeLayout),
		"user_agent", s.UserAgent,
		"client_ip", s.ClientIP,
	)
	if err != nil {
		return err
	}

	expire := s.ExpiresAt.Sub(s.CreatedAt)
	fmt.Println(expire, int(expire), expire.Seconds())
	_, err = conn.Do("EXPIRE", key, int64(expire.Seconds()))
	return err
}

// Save 保存
func (s *Session) Get(pool *redis.Pool, key string) error {
	conn := pool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", key))

	if err != nil {
		return err
	}

	// 判断是否是空
	if len(values) == 0 {
		return ErrSessionNotFound
	}

	// 扫描到结构体
	err = redis.ScanStruct(values, s)
	if err != nil {
		return err
	}

	t1, err := redis.String(conn.Do("HGET", key, "expires_at"))
	if err != nil {
		return err
	}

	t2, err := redis.String(conn.Do("HGET", key, "created_at"))
	if err != nil {
		return err
	}

	s.ExpiresAt, err = time.Parse(timeLayout, t1)
	if err != nil {
		return err
	}

	s.CreatedAt, err = time.Parse(timeLayout, t2)
	if err != nil {
		return err
	}

	return nil
}

// Del 删除
func (s *Session) Del(pool *redis.Pool, key string) error {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}
