package entity

import "gorm.io/gorm"

type Inventory struct {
	gorm.Model
	Stock     uint   `gorm:"not null" json:"stock"`
	Location  string `gorm:"not null" json:"location"`
	ProductId uint   `json:"productId"`
}

type CreateInventoryRequest struct {
	Stock     uint   `json:"stock" binding:"required"`
	Location  string `json:"location" binding:"required"`
	ProductId uint   `json:"productId" binding:"required"`
}

type UpdateInventoryRequest struct {
	Stock uint `json:"stock" binding:"min=0"`
}

type InventoryResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type InventoryDataResponse struct {
	ID        uint   `json:"id"`
	Stock     uint   `json:"stock"`
	Location  string `json:"location"`
	ProductId uint   `json:"productId"`
}
