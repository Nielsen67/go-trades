package repository

import (
	"go-trades/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx *gin.Context, user *entity.User) error
	FindAll(ctx *gin.Context) ([]entity.User, error)
	FindByUsername(ctx *gin.Context, username string) (*entity.User, error)
	FindById(ctx *gin.Context, id uint) (*entity.User, error)
	FindByPhoneNumber(ctx *gin.Context, phoneNumber string) (*entity.User, error)
	FindByEmail(ctx *gin.Context, email string) (*entity.User, error)
	Update(ctx *gin.Context, user *entity.User) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) FindAll(ctx *gin.Context) ([]entity.User, error) {
	var users []entity.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *userRepository) FindByUsername(ctx *gin.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByPhoneNumber(ctx *gin.Context, phonenumber string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where("phonenumber = ?", phonenumber).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx *gin.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindById(ctx *gin.Context, id uint) (*entity.User, error) {
	var user entity.User
	err := r.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx *gin.Context, user *entity.User) error {
	return r.DB.Create(user).Error
}

func (r *userRepository) Update(ctx *gin.Context, user *entity.User) error {
	return r.DB.Save(user).Error
}
