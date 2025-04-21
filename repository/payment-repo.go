package repository

import (
	"go-trades/entity"
	"go-trades/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type paymentRepository struct {
	DB *gorm.DB
}

type PaymentRepository interface {
	FindAll(ctx *gin.Context, page, size int) ([]entity.Payment, int64, error)
	CreatePayment(ctx *gin.Context, payment *entity.Payment) error
	FindAllByUserId(ctx *gin.Context, userId uint, page, size int) ([]entity.Payment, int64, error)
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{
		DB: db,
	}
}

func (r *paymentRepository) FindAll(ctx *gin.Context, page, size int) ([]entity.Payment, int64, error) {
	var result []entity.Payment
	var total int64

	if err := r.DB.Model(&entity.Payment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.DB.Offset(offset).Limit(size).Find(&result).Error
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *paymentRepository) FindAllByUserId(ctx *gin.Context, userId uint, page, size int) ([]entity.Payment, int64, error) {
	var result []entity.Payment
	var total int64

	if err := r.DB.Model(&entity.Payment{}).
		Joins("JOIN orders ON orders.id = payments.order_id").
		Where("orders.user_id = ?", userId).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.DB.Model(&entity.Payment{}).
		Joins("JOIN orders ON orders.id = payments.order_id").
		Where("orders.user_id = ?", userId).
		Offset(offset).
		Limit(size).
		Find(&result).Error
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *paymentRepository) CreatePayment(ctx *gin.Context, payment *entity.Payment) error {
	db := utils.GetTx(ctx, r.DB)
	return db.Create(payment).Error
}
