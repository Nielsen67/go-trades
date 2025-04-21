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
	FindAll(ctx *gin.Context, page, size int) ([]entity.Product, int64, error)
	FindAllWithStock(ctx *gin.Context, page, size int) ([]entity.ProductDataResponse, int64, error)
	FindByCategoryId(ctx *gin.Context, page, size int, id uint) ([]entity.Product, int64, error)
	FindByCategoryIdWithStock(ctx *gin.Context, page, size int, id uint) ([]entity.ProductDataResponse, int64, error)
	FindById(ctx *gin.Context, id uint) (*entity.Product, error)
	FindByIdWithStock(ctx *gin.Context, id uint) (*entity.ProductDataResponse, error)
	FindByName(ctx *gin.Context, name string) (*entity.Product, error)
	CreateProduct(ctx *gin.Context, product *entity.Product) error
	UpdateProduct(ctx *gin.Context, product *entity.Product) error
	DeleteProduct(ctx *gin.Context, id uint) error
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		DB: db,
	}
}

func (r *productRepository) FindAll(ctx *gin.Context, page, size int) ([]entity.Product, int64, error) {
	var result []entity.Product
	var total int64

	if err := r.DB.Model(&entity.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.DB.Offset(offset).Limit(size).Find(&result).Error
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *productRepository) FindAllWithStock(ctx *gin.Context, page, size int) ([]entity.ProductDataResponse, int64, error) {
	var result []entity.ProductDataResponse
	var total int64

	if err := r.DB.Model(&entity.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size

	err := r.DB.Model(&entity.Product{}).
		Select("products.*, COALESCE(SUM(inventories.stock), 0) as stock").
		Joins("LEFT JOIN inventories ON inventories.product_id = products.id").
		Group("products.id").
		Where("inventories.deleted_at IS NULL").
		Offset(offset).
		Limit(size).
		Find(&result).Error
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *productRepository) FindByCategoryId(ctx *gin.Context, page, size int, id uint) ([]entity.Product, int64, error) {
	var result []entity.Product
	var total int64

	if err := r.DB.Model(&entity.Product{}).Where("category_id = ?", id).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size

	err := r.DB.Offset(offset).Limit(size).Where("category_id = ?", id).Find(&result).Error
	if err != nil {
		return nil, 0, err
	}
	return result, total, nil
}

func (r *productRepository) FindByCategoryIdWithStock(ctx *gin.Context, page, size int, id uint) ([]entity.ProductDataResponse, int64, error) {
	var result []entity.ProductDataResponse
	var total int64

	if err := r.DB.Model(&entity.Product{}).Where("category_id = ?", id).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size

	err := r.DB.Model(&entity.Product{}).
		Select("products.*, COALESCE(SUM(inventories.stock), 0) as stock").
		Joins("LEFT JOIN inventories ON inventories.product_id = products.id").
		Where("products.category_id = ? AND inventories.deleted_at IS NULL", id).
		Group("products.id").
		Offset(offset).
		Limit(size).
		Find(&result).Error
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
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

func (r *productRepository) FindByIdWithStock(ctx *gin.Context, id uint) (*entity.ProductDataResponse, error) {
	var result entity.ProductDataResponse
	db := utils.GetTx(ctx, r.DB)

	err := db.Model(&entity.Product{}).
		Select("products.*, COALESCE(SUM(inventories.stock), 0) as stock").
		Joins("LEFT JOIN inventories ON inventories.product_id = products.id").
		Where("products.id = ? AND inventories.deleted_at IS NULL", id).
		Group("products.id").
		First(&result).Error
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
