package repository

import (
	"errors"
	"go-trades/entity"

	"gorm.io/gorm"
)

type productRepository struct {
	DB *gorm.DB
}

type ProductRepository interface {
	FindAll() ([]entity.Product, error)
	FindByCategoryId(id uint) ([]entity.Product, error)
	FindById(id uint) (*entity.Product, error)
	FindByName(name string) (*entity.Product, error)
	CreateProduct(product *entity.Product) error
	UpdateProduct(product *entity.Product) error
	DeleteProduct(id uint) error
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		DB: db,
	}
}

func (r *productRepository) FindAll() ([]entity.Product, error) {
	var result []entity.Product
	err := r.DB.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *productRepository) FindByCategoryId(id uint) ([]entity.Product, error) {
	var result []entity.Product
	err := r.DB.Where("category_id = ?", id).Find(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *productRepository) FindById(id uint) (*entity.Product, error) {
	var result entity.Product
	err := r.DB.Where("id = ?", id).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *productRepository) FindByName(name string) (*entity.Product, error) {
	var result entity.Product
	err := r.DB.Where("name = ?", name).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *productRepository) CreateProduct(product *entity.Product) error {
	return r.DB.Create(product).Error
}

func (r *productRepository) UpdateProduct(product *entity.Product) error {
	return r.DB.Save(product).Error
}

func (r *productRepository) DeleteProduct(id uint) error {
	var product entity.Product
	if err := r.DB.First(&product, id).Error; err != nil {
		return err
	}
	return r.DB.Delete(&product).Error
}
