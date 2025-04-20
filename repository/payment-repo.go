package repository

import (
	"errors"
	"go-trades/entity"
	"go-trades/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type paymentRepository struct {
	DB *gorm.DB
}

type PaymentRepository interface {
	FindAll(ctx *gin.Context) ([]entity.Payment, error)
	FindById(ctx *gin.Context, id uint) (*entity.Payment, error)
	FindByOrderId(ctx *gin.Context, id uint) (*entity.Payment, error)
	CreatePayment(ctx *gin.Context, payment *entity.Payment) error
	resolveDB(tx *gorm.DB) *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{
		DB: db,
	}
}

func (r *paymentRepository) resolveDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.DB
}

func (r *paymentRepository) FindAll(ctx *gin.Context) ([]entity.Payment, error) {
	var result []entity.Payment
	err := r.DB.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *paymentRepository) FindById(ctx *gin.Context, id uint) (*entity.Payment, error) {
	var result entity.Payment
	err := r.DB.Where("id = ?", id).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *paymentRepository) FindByOrderId(ctx *gin.Context, id uint) (*entity.Payment, error) {
	var result entity.Payment
	err := r.DB.Where("order_id = ?", id).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *paymentRepository) CreatePayment(ctx *gin.Context, payment *entity.Payment) error {
	db := utils.GetTx(ctx, r.DB)
	return db.Create(payment).Error
}
