package dto

import (
	"encoding/json"
	"log"
)

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

type LoginRequest struct {
	Email     string `json:"email"  binding:"required,email"`
	Password  string `json:"password"  binding:"required,min=6,max=32"`
	UserAgent string `binding:"-"`
	ClientIP  string `binding:"-"`
}

type LoginResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
