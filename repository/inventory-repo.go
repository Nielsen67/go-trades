package repository

import (
	"errors"
	"go-trades/entity"
	"go-trades/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type inventoryRepository struct {
	DB *gorm.DB
}

type InventoryRepository interface {
	FindAll(ctx *gin.Context, page, size int) ([]entity.Inventory, int64, error)
	FindById(ctx *gin.Context, id uint) (*entity.Inventory, error)
	FindFirstByProductId(ctx *gin.Context, id uint) (*entity.Inventory, error)
	FindByName(ctx *gin.Context, name string) (*entity.Inventory, error)
	FindByCode(ctx *gin.Context, code string) (*entity.Inventory, error)
	CreateInventory(ctx *gin.Context, inventory *entity.Inventory) error
	UpdateInventory(ctx *gin.Context, inventory *entity.Inventory) error
	UpdateInventoryForOrder(ctx *gin.Context, inventory *entity.Inventory, qty uint, action string) error
	DeleteInventory(ctx *gin.Context, id uint) error
	resolveDB(tx *gorm.DB) *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{
		DB: db,
	}
}

func (r *inventoryRepository) resolveDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.DB
}

func (r *inventoryRepository) FindAll(ctx *gin.Context, page, size int) ([]entity.Inventory, int64, error) {
	var result []entity.Inventory
	var total int64

	if err := r.DB.Model(&entity.Inventory{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.DB.Offset(offset).Limit(size).Find(&result).Error
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *inventoryRepository) FindById(ctx *gin.Context, id uint) (*entity.Inventory, error) {
	var result entity.Inventory
	err := r.DB.Where("id = ?", id).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *inventoryRepository) FindByCode(ctx *gin.Context, code string) (*entity.Inventory, error) {
	var result entity.Inventory
	err := r.DB.Where("code = ?", code).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *inventoryRepository) FindByName(ctx *gin.Context, name string) (*entity.Inventory, error) {
	var result entity.Inventory
	err := r.DB.Where("name = ?", name).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *inventoryRepository) FindFirstByProductId(ctx *gin.Context, id uint) (*entity.Inventory, error) {
	var result entity.Inventory
	db := utils.GetTx(ctx, r.DB)
	err := db.Where("product_id = ?", id).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *inventoryRepository) CreateInventory(ctx *gin.Context, inventory *entity.Inventory) error {
	return r.DB.Create(inventory).Error
}

func (r *inventoryRepository) UpdateInventory(ctx *gin.Context, inventory *entity.Inventory) error {
	return r.DB.Save(inventory).Error
}

func (r *inventoryRepository) UpdateInventoryForOrder(ctx *gin.Context, inventory *entity.Inventory, qty uint, action string) error {
	db := utils.GetTx(ctx, r.DB)

	switch action {
	case "create":
		if err := db.Model(&inventory).Where("id = ?", inventory.ID).Update("stock", gorm.Expr("stock - ?", qty)).Error; err != nil {
			return err
		}
	case "cancel":
		if err := db.Model(&inventory).Where("id = ?", inventory.ID).Update("stock", gorm.Expr("stock + ?", qty)).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *inventoryRepository) DeleteInventory(ctx *gin.Context, id uint) error {
	var inventory entity.Inventory
	if err := r.DB.First(&inventory, id).Error; err != nil {
		return err
	}
	return r.DB.Delete(&inventory).Error
}
