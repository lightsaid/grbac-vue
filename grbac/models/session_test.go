package models

import (
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/require"
)

func TestRedis(t *testing.T) {
	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 5 * time.Minute,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6399")
		},
	}

	var key = "test_session"

	ss := Session{
		TID:          "1",
		UID:          1,
		RefreshToken: "aabbbcc",
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(time.Second * 20),
	}

	err := ss.Save(pool, key)
	require.NoError(t, err)

	var s3 Session
	err = s3.Get(pool, key)
	require.NoError(t, err)
	require.Equal(t, ss.TID, s3.TID)
	require.Equal(t, ss.UID, s3.UID)
	require.Equal(t, ss.RefreshToken, s3.RefreshToken)
	require.EqualValues(t, ss.CreatedAt.Format(timeLayout), s3.CreatedAt.Format(timeLayout))
	require.EqualValues(t, ss.ExpiresAt.Format(timeLayout), s3.ExpiresAt.Format(timeLayout))
}
