package data

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

type SessionModel struct {
	Pool *redis.Pool
}

// Save 保存
func (m *SessionModel) Save(key string, sess Session) error {
	conn := m.Pool.Get()
	defer conn.Close()

	_, err := conn.Do(
		"HMSET", key,
		"tid", sess.TID,
		"uid", sess.UID,
		"refresh_token", sess.RefreshToken,
		"expires_at", sess.ExpiresAt.Format(timeLayout),
		"created_at", sess.CreatedAt.Format(timeLayout),
		"user_agent", sess.UserAgent,
		"client_ip", sess.ClientIP,
	)
	if err != nil {
		return err
	}

	expire := sess.ExpiresAt.Sub(sess.CreatedAt)
	fmt.Println(expire, int(expire), expire.Seconds())
	_, err = conn.Do("EXPIRE", key, int64(expire.Seconds()))
	return err
}

// Save 保存
func (m *SessionModel) Get(key string) (*Session, error) {
	conn := m.Pool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", key))

	if err != nil {
		return nil, err
	}

	// 判断是否是空
	if len(values) == 0 {
		return nil, ErrSessionNotFound
	}

	var sess Session

	// 扫描到结构体
	err = redis.ScanStruct(values, sess)
	if err != nil {
		return nil, err
	}

	t1, err := redis.String(conn.Do("HGET", key, "expires_at"))
	if err != nil {
		return nil, err
	}

	t2, err := redis.String(conn.Do("HGET", key, "created_at"))
	if err != nil {
		return nil, err
	}

	sess.ExpiresAt, err = time.Parse(timeLayout, t1)
	if err != nil {
		return nil, err
	}

	sess.CreatedAt, err = time.Parse(timeLayout, t2)
	if err != nil {
		return nil, err
	}

	return &sess, nil
}

// Del 删除
func (s *SessionModel) Del(pool *redis.Pool, key string) error {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}
