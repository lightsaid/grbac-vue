package data

import "gorm.io/gorm"

// Permission 权限表结构
type Permission struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Code        string `json:"code" gorm:"type:varchar(64);uniqueIndex;not null"`
	Name        string `json:"name" gorm:"type:varchar(64);not null"`
	Description string `json:"description" gorm:"type:varchar(256)"`
	*BaseModel
}

type PermissionModel struct {
	DB *gorm.DB
}
