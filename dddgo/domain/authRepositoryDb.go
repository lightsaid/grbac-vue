package domain

import (
	"errors"
	"time"

	"github.com/lightsaid/grbac-vue/dddgo/entity"
	"gorm.io/gorm"
)

var ErrUserIsActived = errors.New("已经激活")

type AuthRepository interface {
	ActivatedAccount(verifyCode string) (user *entity.User, err error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) ActivatedAccount(verifyCode string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("verify_code = ?", verifyCode).First(&user).Error
	if err != nil {
		return nil, err
	}

	if user.ActivatedAt != nil {
		return &user, ErrUserIsActived
	}

	err = r.db.Model(&user).Where("verify_code = ?", verifyCode).Update("activated_at", time.Now()).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
