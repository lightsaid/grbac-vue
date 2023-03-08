package entity

import "github.com/lightsaid/grbac-vue/dddgo/dto"

// Role 角色表结构
type Role struct {
	ID          uint          `gorm:"primaryKey"`
	Code        string        `gorm:"type:varchar(64);uniqueIndex;not null"`
	Name        string        `gorm:"type:varchar(64);not null"`
	Description string        `gorm:"type:varchar(256)"`
	Permissions []*Permission `gorm:"many2many:role_permissions"`
	*Base
}

func (r *Role) ToDto() *dto.RoleResponse {
	var permissions []*dto.PermissionResponse
	for _, p := range r.Permissions {
		permissions = append(permissions, p.ToDto())
	}
	return &dto.RoleResponse{
		ID:          r.ID,
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		Permissions: permissions,
	}
}
