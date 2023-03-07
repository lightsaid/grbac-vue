package models

import (
	"bytes"
	"html/template"
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
	VerifyCode  string     `json:"-"`
	Roles       []*Role    `gorm:"many2many:user_roles"`
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
