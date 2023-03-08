package domain

import (
	"github.com/lightsaid/grbac-vue/dddgo/entity"
	"gorm.io/gorm"
)

// UserRepository 定义 UserRepository 要实现的接口，提供给service
type UserRepository interface {
	Create(user *entity.User) error
	FindByID(id uint) (user *entity.User, err error)
	FindByEmail(email string) (user *entity.User, err error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Roles").Where("email = ?", email).First(&user).Error
	return &user, err
}
