package data

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

// BaseModel 基础模型，公共模型
type BaseModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Models struct {
	User       UserModel
	Role       RoleModel
	Permission PermissionModel
	Session    SessionModel
}

func NewModels(db *gorm.DB, pool *redis.Pool) Models {
	return Models{
		User:       UserModel{DB: db},
		Role:       RoleModel{DB: db},
		Permission: PermissionModel{DB: db},
		Session:    SessionModel{Pool: pool},
	}
}
