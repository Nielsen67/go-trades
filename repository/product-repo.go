package repository

import (
	"errors"
	"go-trades/entity"
	"go-trades/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type productRepository struct {
	DB *gorm.DB
}

type ProductRepository interface {
	FindAll(ctx *gin.Context) ([]entity.Product, error)
	FindByCategoryId(ctx *gin.Context, id uint) ([]entity.Product, error)
	FindById(ctx *gin.Context, id uint) (*entity.Product, error)
	FindByName(ctx *gin.Context, name string) (*entity.Product, error)
	CreateProduct(ctx *gin.Context, product *entity.Product) error
	UpdateProduct(ctx *gin.Context, product *entity.Product) error
	DeleteProduct(ctx *gin.Context, id uint) error
	resolveDB(tx *gorm.DB) *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		DB: db,
	}
}

func (r *productRepository) resolveDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.DB
}

func (r *productRepository) FindAll(ctx *gin.Context) ([]entity.Product, error) {
	var result []entity.Product
	err := r.DB.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *productRepository) FindByCategoryId(ctx *gin.Context, id uint) ([]entity.Product, error) {
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

func (r *productRepository) FindById(ctx *gin.Context, id uint) (*entity.Product, error) {
	var result entity.Product
	db := utils.GetTx(ctx, r.DB)

	err := db.Where("id = ?", id).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *productRepository) FindByName(ctx *gin.Context, name string) (*entity.Product, error) {
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

func (r *productRepository) CreateProduct(ctx *gin.Context, product *entity.Product) error {
	return r.DB.Create(product).Error
}

func (r *productRepository) UpdateProduct(ctx *gin.Context, product *entity.Product) error {
	return r.DB.Save(product).Error
}

func (r *productRepository) DeleteProduct(ctx *gin.Context, id uint) error {
	var product entity.Product
	if err := r.DB.First(&product, id).Error; err != nil {
		return err
	}
	return r.DB.Delete(&product).Error
}
