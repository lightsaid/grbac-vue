package data

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表结构
type User struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"type:varchar(16);not null"`
	Email       string     `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password    string     `json:"-" gorm:"type:varchar(64);not null;"`
	Avatar      *string    `json:"avatar" gorm:"type:varchar(255)"`
	ActivatedAt *time.Time `json:"activated_at"`
	VerifyCode  *string    `json:"-" gorm:"type:varchar(255)"`
	Roles       []*Role    `jsoon:"roles" gorm:"many2many:user_roles"`
	*BaseModel
}

type UserModel struct {
	DB *gorm.DB
}
