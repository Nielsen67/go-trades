package service

import (
	"errors"
	"go-trades/config"
	"go-trades/entity"
	"go-trades/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type userService struct {
	Repository repository.UserRepository
}

type UserService interface {
	GetUserById(ctx *gin.Context, userId uint) (*entity.UserDataResponse, error)
	Register(ctx *gin.Context, req *entity.UserRegisterRequest) (*entity.UserDataResponse, error)
	Login(ctx *gin.Context, req *entity.UserLoginRequest) (*entity.UserLoginResponse, error)
	ChangePassword(ctx *gin.Context, userId uint, req *entity.UserChangePassword) (*entity.UserChangePasswordResponse, error)
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{
		Repository: r,
	}
}

func (s *userService) GetUserById(ctx *gin.Context, userId uint) (*entity.UserDataResponse, error) {
	user, err := s.Repository.FindById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &entity.UserDataResponse{
		Id:          int(user.ID),
		Username:    user.Username,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Dob:         user.Dob,
		Address:     user.Address,
		Email:       user.Email,
		Phonenumber: user.Phonenumber,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (s *userService) Register(ctx *gin.Context, req *entity.UserRegisterRequest) (*entity.UserDataResponse, error) {

	var user *entity.User
	var err error

	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	user, err = s.Repository.FindByUsername(ctx, req.Username)
	if user != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("username already taken")
	}

	user, err = s.Repository.FindByEmail(ctx, req.Email)
	if user != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("email already taken")
	}

	user, err = s.Repository.FindByPhoneNumber(ctx, req.Phonenumber)
	if user != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("phone number already taken")
	}

	user = &entity.User{
		Username:    req.Username,
		Firstname:   req.Firstname,
		Lastname:    req.Lastname,
		Dob:         dob,
		Address:     req.Address,
		Email:       req.Email,
		Phonenumber: req.Phonenumber,
		Role:        entity.Customer,
	}

	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.Repository.Create(ctx, user); err != nil {
		return nil, err
	}

	return &entity.UserDataResponse{
		Id:          int(user.ID),
		Username:    user.Username,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Dob:         user.Dob,
		Address:     user.Address,
		Email:       user.Email,
		Phonenumber: user.Phonenumber,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (s *userService) Login(ctx *gin.Context, req *entity.UserLoginRequest) (*entity.UserLoginResponse, error) {
	user, err := s.Repository.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("username/password is invalid")
	}

	if err := user.CheckPassword(req.Password); err != nil {
		return nil, errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"role":   user.Role,
		"exp":    time.Now().Add(config.GetJWTExpirationDuration()).Unix(),
	})

	tokenString, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		return nil, err
	}

	return &entity.UserLoginResponse{Token: tokenString}, nil
}

func (s *userService) ChangePassword(ctx *gin.Context, userId uint, req *entity.UserChangePassword) (*entity.UserChangePasswordResponse, error) {
	user, err := s.Repository.FindById(ctx, userId)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := user.CheckPassword(req.OldPassword); err != nil {
		return nil, errors.New("invalid old password")
	}

	if err := user.HashPassword(req.NewPassword); err != nil {
		return nil, err
	}

	if err := s.Repository.Update(ctx, user); err != nil {
		return nil, err
	}

	return &entity.UserChangePasswordResponse{Message: "Password changed successfully"}, nil
}
