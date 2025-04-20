package repository

import (
	"errors"
	"go-trades/entity"
	"go-trades/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type orderRepository struct {
	DB *gorm.DB
}

type OrderRepository interface {
	FindAll(ctx *gin.Context) ([]entity.Order, error)
	FindById(ctx *gin.Context, id uint) (*entity.Order, error)
	FindByStatus(ctx *gin.Context, status uint) ([]entity.Order, error)
	CreateOrder(ctx *gin.Context, order *entity.Order) error
	UpdateOrder(ctx *gin.Context, order *entity.Order) error
	DeleteOrder(ctx *gin.Context, id uint) error
	resolveDB(tx *gorm.DB) *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		DB: db,
	}
}

func (r *orderRepository) resolveDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.DB
}

func (r *orderRepository) FindAll(ctx *gin.Context) ([]entity.Order, error) {
	var result []entity.Order
	err := r.DB.Preload("OrderDetails").Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *orderRepository) FindById(ctx *gin.Context, id uint) (*entity.Order, error) {
	var result entity.Order
	err := r.DB.Preload("OrderDetails").Where("id = ?", id).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *orderRepository) FindByStatus(ctx *gin.Context, status uint) ([]entity.Order, error) {
	var result []entity.Order
	err := r.DB.Preload("OrderDetails").Where("status = ?", status).Find(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *orderRepository) CreateOrder(ctx *gin.Context, order *entity.Order) error {
	db := utils.GetTx(ctx, r.DB)
	return db.Create(order).Error
}

func (r *orderRepository) UpdateOrder(ctx *gin.Context, order *entity.Order) error {
	db := utils.GetTx(ctx, r.DB)
	return db.Save(order).Error
}

func (r *orderRepository) DeleteOrder(ctx *gin.Context, id uint) error {
	var order entity.Order
	if err := r.DB.First(&order, id).Error; err != nil {
		return err
	}
	return r.DB.Delete(&order).Error
}
