package repository

import (
	"go-trades/entity"
	"go-trades/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type productImageRepository struct {
	DB *gorm.DB
}

type ProductImageRepository interface {
	CreateProductImage(ctx *gin.Context, productImage *entity.ProductImage) error
	FindAllByProductId(ctx *gin.Context, productId uint) ([]entity.ProductImage, error)
}

func NewProductImageRepository(db *gorm.DB) ProductImageRepository {
	return &productImageRepository{
		DB: db,
	}
}

func (r *productImageRepository) CreateProductImage(ctx *gin.Context, productImage *entity.ProductImage) error {
	db := utils.GetTx(ctx, r.DB)
	return db.Create(productImage).Error
}

func (r *productImageRepository) FindAllByProductId(ctx *gin.Context, productId uint) ([]entity.ProductImage, error) {
	var productImages []entity.ProductImage
	db := utils.GetTx(ctx, r.DB)
	err := db.Where("product_id = ?", productId).Find(&productImages).Error
	if err != nil {
		return nil, err
	}
	return productImages, nil
}
