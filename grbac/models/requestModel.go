package models

import (
	"encoding/json"
	"log"
)

type RegisterRequest struct {
	Name            string `json:"name" zh:"用户名" binding:"required,min=2,max=32"`
	Email           string `json:"email" zh:"邮箱地址" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6,max=32"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password"  binding:"required,min=6,max=32"`
}

type LoginResponse struct {
	User
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ActivateUserRequest struct {
	VerifyCode string `uri:"verifyCode" binding:"required"`
}

type RegisterMailerPaylod struct {
	Email   string `json:"email"`
	Content string `json:"content"`
}

func (r *RegisterMailerPaylod) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
	}
	return string(b)
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateProfileRequest struct {
	Name   string `json:"name" binding:"required"`
	Avatar string `json:"avatar" binding:"required"`
}

// 分页
type Pagination struct {
	Page int `form:"page" binding:"required"`
	Size int `form:"size" binding:"required"`
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Size
}

func (p *Pagination) Limit() int {
	return p.Size
}

// PermissionRequest Role的CRUD入参
type PermissionRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"-"`
}

// RoleRequest Role的CRUD入参
type RoleRequest struct {
	Code          string `json:"code" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description" binding:"-"`
	PermissionIds []uint `json:"permissions" binding:"-"`
}

// RequestUri   v1/api/xxx/:id
type RequestUri struct {
	ID uint `uri:"id"`
}
