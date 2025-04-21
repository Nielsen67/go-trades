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
	FindAll(ctx *gin.Context, page, size int) ([]entity.Order, int64, error)
	FindById(ctx *gin.Context, id uint) (*entity.Order, error)
	FindByStatus(ctx *gin.Context, page, size int, status uint) ([]entity.Order, int64, error)
	FindAllByUserId(ctx *gin.Context, userId uint, page, size int) ([]entity.Order, int64, error)
	FindByUserIdWithId(ctx *gin.Context, userId uint, id uint) (*entity.Order, error)
	FindAllByUserIdWithStatus(ctx *gin.Context, userId uint, page, size int, status uint) ([]entity.Order, int64, error)
	CreateOrder(ctx *gin.Context, order *entity.Order) error
	UpdateOrder(ctx *gin.Context, order *entity.Order) error
	DeleteOrder(ctx *gin.Context, id uint) error
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		DB: db,
	}
}

func (r *orderRepository) FindAll(ctx *gin.Context, page, size int) ([]entity.Order, int64, error) {
	var result []entity.Order
	var total int64

	if err := r.DB.Model(&entity.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.DB.Preload("OrderDetails").Offset(offset).Limit(size).Find(&result).Error
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
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

func (r *orderRepository) FindByStatus(ctx *gin.Context, page, size int, status uint) ([]entity.Order, int64, error) {
	var result []entity.Order
	var total int64

	if err := r.DB.Model(&entity.Order{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size

	err := r.DB.Preload("OrderDetails").Offset(offset).Limit(size).Where("status = ?", status).Find(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	return result, total, nil
}

func (r *orderRepository) FindAllByUserId(ctx *gin.Context, userId uint, page, size int) ([]entity.Order, int64, error) {
	var result []entity.Order
	var total int64

	if err := r.DB.Model(&entity.Order{}).Where("user_id = ?", userId).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.DB.Preload("OrderDetails").Where("user_id = ?", userId).Offset(offset).Limit(size).Find(&result).Error
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *orderRepository) FindByUserIdWithId(ctx *gin.Context, userId uint, id uint) (*entity.Order, error) {
	var result entity.Order

	err := r.DB.Preload("OrderDetails").Where("user_id = ? AND id = ?", userId, id).First(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *orderRepository) FindAllByUserIdWithStatus(ctx *gin.Context, userId uint, page, size int, status uint) ([]entity.Order, int64, error) {
	var result []entity.Order
	var total int64

	if err := r.DB.Model(&entity.Order{}).Where("user_id = ? AND status = ?", userId, status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size

	err := r.DB.Preload("OrderDetails").Where("user_id = ? AND status = ?", userId, status).Offset(offset).Limit(size).Find(&result).Error
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
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
	db := utils.GetTx(ctx, r.DB)

	var order entity.Order
	if err := db.First(&order, id).Error; err != nil {
		return err
	}

	if err := db.Where("order_id = ?", id).Delete(&entity.OrderDetail{}).Error; err != nil {
		return err
	}

	return db.Delete(&order).Error
}
