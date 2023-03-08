package dto

import (
	"time"
)

type UserResponse struct {
	ID        uint            `json:"id"`
	Email     string          `json:"email"`
	Avatar    *string         `json:"avatar"`
	Roles     []*RoleResponse `json:"roles"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type NewUserRequest struct {
	Name            string `json:"name" zh:"用户名" binding:"required,min=2,max=32"`
	Email           string `json:"email" zh:"邮箱地址" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6,max=32"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
}
