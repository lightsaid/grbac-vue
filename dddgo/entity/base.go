package entity

import (
	"time"

	"gorm.io/gorm"
)

// Base 基础模型，公共模型
type Base struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
