package repository

import (
	"errors"
	"go-trades/entity"

	"gorm.io/gorm"
)

type inventoryRepository struct {
	DB *gorm.DB
}

type InventoryRepository interface {
	FindAll() ([]entity.Inventory, error)
	FindById(id uint) (*entity.Inventory, error)
	FindByName(name string) (*entity.Inventory, error)
	FindByCode(code string) (*entity.Inventory, error)
	CreateInventory(inventory *entity.Inventory) error
	UpdateInventory(inventory *entity.Inventory) error
	DeleteInventory(id uint) error
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{
		DB: db,
	}
}

func (r *inventoryRepository) FindAll() ([]entity.Inventory, error) {
	var result []entity.Inventory
	err := r.DB.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *inventoryRepository) FindById(id uint) (*entity.Inventory, error) {
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

func (r *inventoryRepository) FindByCode(code string) (*entity.Inventory, error) {
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

func (r *inventoryRepository) FindByName(name string) (*entity.Inventory, error) {
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

func (r *inventoryRepository) CreateInventory(inventory *entity.Inventory) error {
	return r.DB.Create(inventory).Error
}

func (r *inventoryRepository) UpdateInventory(inventory *entity.Inventory) error {
	return r.DB.Save(inventory).Error
}

func (r *inventoryRepository) DeleteInventory(id uint) error {
	var inventory entity.Inventory
	if err := r.DB.First(&inventory, id).Error; err != nil {
		return err
	}
	return r.DB.Delete(&inventory).Error
}
