package entity

import "github.com/lightsaid/grbac-vue/dddgo/dto"

// Permission 权限表结构
type Permission struct {
	ID          uint    `gorm:"primaryKey"`
	Code        string  `gorm:"type:varchar(64);uniqueIndex;not null"`
	Name        string  `gorm:"type:varchar(64);not null"`
	Description *string `gorm:"type:varchar(256)"`
	*Base
}

func (p Permission) ToDto() *dto.PermissionResponse {
	return &dto.PermissionResponse{
		ID:          p.ID,
		Code:        p.Code,
		Name:        p.Name,
		Description: *p.Description,
	}
}
