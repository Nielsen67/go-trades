package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role string

const (
	Admin    Role = "admin"
	Customer Role = "customer"
)

type User struct {
	gorm.Model
	Username    string    `gorm:"unique;not null;size:50" json:"username"`
	Password    string    `gorm:"not null" json:"password"`
	Firstname   string    `gorm:"not null" json:"firstname"`
	Lastname    string    `gorm:"not null" json:"lastname"`
	Dob         time.Time `gorm:"not null;type:date" json:"dob"`
	Address     string    `gorm:"not null" json:"address"`
	Email       string    `gorm:"unique;not null;size:255" json:"email"`
	Phonenumber string    `gorm:"unique;not null" json:"phoneNumber"`
	Role        Role      `gorm:"not null;type:enum('admin', 'customer');default:customer" json:"role"`
	Token       string    `gorm:"default:null" json:"access_token"`
	Orders      []Order   `gorm:"foreignKey:UserId"`
}

type UserRegisterRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Firstname   string `json:"firstname" binding:"required"`
	Lastname    string `json:"lastname" binding:"required"`
	Dob         string `json:"dob" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Phonenumber string `json:"phoneNumber" binding:"required"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type UserDataResponse struct {
	Id          int       `json:"id"`
	Username    string    `json:"username"`
	Firstname   string    `json:"firstname"`
	Lastname    string    `json:"lastname"`
	Dob         time.Time `json:"dob"`
	Address     string    `json:"address"`
	Email       string    `json:"email"`
	Phonenumber string    `json:"phoneNumber"`
	Role        Role      `json:"role"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UserChangePassword struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UserChangePasswordResponse struct {
	Message string `json:"message"`
}

type UserLoginResponse struct {
	Token string `json:"access_token"`
}

func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
