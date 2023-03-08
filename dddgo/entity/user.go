package entity

import (
	"bytes"
	"html/template"
	"time"

	"github.com/lightsaid/grbac-vue/dddgo/dto"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	ActivateEmailTemplatePath = "./views/activate.mailer.tmpl"
)

// User 用户表结构
type User struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"type:varchar(16);not null"`
	Email       string  `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password    string  `gorm:"type:varchar(64);not null;"`
	Avatar      *string `gorm:"type:varchar(255)"`
	ActivatedAt *time.Time
	VerifyCode  *string `gorm:"type:varchar(255);comment '注册时激活账号发送邮件的签名'"`
	Roles       []*Role `gorm:"many2many:user_roles"`
	*Base
}

func (u *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	u.Password = string(hashedPassword)
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
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

func (user *User) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(&User{}).Count(&total)

	return total
}

func (u *User) ToDto() *dto.UserResponse {
	var roles []*dto.RoleResponse
	for _, r := range u.Roles {
		roles = append(roles, r.ToDto())
	}
	return &dto.UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Avatar:    u.Avatar,
		Roles:     roles,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
