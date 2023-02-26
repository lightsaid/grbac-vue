package models

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	ActivateEmailTemplatePath = "./views/activate.mailer.tmpl"
)

type User struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"type:varchar(16);not null"`
	Email       string     `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password    string     `json:"-" gorm:"type:varchar(64);not null;"`
	Avatar      string     `json:"avatar" gorm:"type:varchar(255)"`
	UserType    uint8      `json:"user_type" gorm:"type:tinyint;not null;default:1;index,comment '1表示普通用户, 2表示管理员'"`
	ActivatedAt *time.Time `json:"activated"`
	VerifyCode  string     `json:"verify_code"`
	*BaseModel
}

func (user *User) SetPassword(passsword string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(passsword), 10)
	if err != nil {
		return err
	}
	user.Password = string(hashedPass)
	return nil
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func (user *User) SetActivateEmailMessage(activateEmailURL string) (string, error) {
	t, err := template.ParseFiles(ActivateEmailTemplatePath)
	if err != nil {
		return "", err
	}

	// 解析邮箱模板
	var buf bytes.Buffer
	err = t.Execute(&buf, struct{ Href string }{Href: activateEmailURL})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type RegisterRequest struct {
	Name            string `json:"name" zh:"用户名" binding:"required,min=2,max=32"`
	Email           string `json:"email" zh:"邮箱地址" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6,max=32"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password"  binding:"required,min=8,max=32"`
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
