package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/lightsaid/grbac-vue/dddgo/entity"
)

var ErrSessionNotFound = errors.New("not found")

const (
	timeLayout = "2006-01-02 15:04:05"
)

type SessionRepository interface {
	Save(s entity.Session, key string) error
	Get(key string) (*entity.Session, error)
	Del(key string) error
}

type sessionRepository struct {
	pool *redis.Pool
}

func NewSessionRepository(pool *redis.Pool) SessionRepository {
	return &sessionRepository{pool: pool}
}

// Save 保存
func (repo *sessionRepository) Save(s entity.Session, key string) error {
	conn := repo.pool.Get()
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
func (repo *sessionRepository) Get(key string) (*entity.Session, error) {
	conn := repo.pool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", key))

	if err != nil {
		return nil, err
	}

	// 判断是否是空
	if len(values) == 0 {
		return nil, ErrSessionNotFound
	}
	var sess entity.Session
	// 扫描到结构体
	err = redis.ScanStruct(values, &sess)
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
func (repo *sessionRepository) Del(key string) error {
	conn := repo.pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}
