package data

import "gorm.io/gorm"

// Role 角色表结构
type Role struct {
	ID          uint          `json:"id" gorm:"primaryKey"`
	Code        string        `json:"code" gorm:"type:varchar(64);uniqueIndex;not null"`
	Name        string        `json:"name" gorm:"type:varchar(64);not null"`
	Description string        `json:"description" gorm:"type:varchar(256)"`
	Permissions []*Permission `json:"permissions" gorm:"many2many:role_permissions"`
	*BaseModel
}

type RoleModel struct {
	DB *gorm.DB
}
