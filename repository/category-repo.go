package repository

import (
	"errors"
	"go-trades/entity"

	"gorm.io/gorm"
)

type categoryRepository struct {
	DB *gorm.DB
}

type CategoryRepository interface {
	FindAll() ([]entity.Category, error)
	FindById(id uint) (*entity.Category, error)
	FindByName(name string) (*entity.Category, error)
	FindByCode(code string) (*entity.Category, error)
	CreateCategory(category *entity.Category) error
	UpdateCategory(category *entity.Category) error
	DeleteCategory(id uint) error
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		DB: db,
	}
}

func (r *categoryRepository) FindAll() ([]entity.Category, error) {
	var result []entity.Category
	err := r.DB.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *categoryRepository) FindById(id uint) (*entity.Category, error) {
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

func (r *categoryRepository) FindByCode(code string) (*entity.Category, error) {
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

func (r *categoryRepository) FindByName(name string) (*entity.Category, error) {
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

func (r *categoryRepository) CreateCategory(category *entity.Category) error {
	return r.DB.Create(category).Error
}

func (r *categoryRepository) UpdateCategory(category *entity.Category) error {
	return r.DB.Save(category).Error
}

func (r *categoryRepository) DeleteCategory(id uint) error {
	var category entity.Category
	if err := r.DB.First(&category, id).Error; err != nil {
		return err
	}
	return r.DB.Delete(&category).Error
}
