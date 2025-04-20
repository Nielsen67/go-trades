package repository

import (
	"errors"
	"go-trades/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type categoryRepository struct {
	DB *gorm.DB
}

type CategoryRepository interface {
	FindAll(ctx *gin.Context) ([]entity.Category, error)
	FindById(ctx *gin.Context, id uint) (*entity.Category, error)
	FindByName(ctx *gin.Context, name string) (*entity.Category, error)
	FindByCode(ctx *gin.Context, code string) (*entity.Category, error)
	CreateCategory(ctx *gin.Context, category *entity.Category) error
	UpdateCategory(ctx *gin.Context, category *entity.Category) error
	DeleteCategory(ctx *gin.Context, id uint) error
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		DB: db,
	}
}

func (r *categoryRepository) FindAll(ctx *gin.Context) ([]entity.Category, error) {
	var result []entity.Category
	err := r.DB.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *categoryRepository) FindById(ctx *gin.Context, id uint) (*entity.Category, error) {
	var result entity.Category
	err := r.DB.Where("id = ?", id).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *categoryRepository) FindByCode(ctx *gin.Context, code string) (*entity.Category, error) {
	var result entity.Category
	err := r.DB.Where("code = ?", code).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *categoryRepository) FindByName(ctx *gin.Context, name string) (*entity.Category, error) {
	var result entity.Category
	err := r.DB.Where("name = ?", name).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *categoryRepository) CreateCategory(ctx *gin.Context, category *entity.Category) error {
	return r.DB.Create(category).Error
}

func (r *categoryRepository) UpdateCategory(ctx *gin.Context, category *entity.Category) error {
	return r.DB.Save(category).Error
}

func (r *categoryRepository) DeleteCategory(ctx *gin.Context, id uint) error {
	var category entity.Category
	if err := r.DB.First(&category, id).Error; err != nil {
		return err
	}
	return r.DB.Delete(&category).Error
}
