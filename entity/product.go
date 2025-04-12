package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name          string         `gorm:"not null;unique" json:"name"`
	Description   string         `json:"description"`
	Price         uint           `gorm:"not null" json:"price"`
	Category      string         `gorm:"not null" json:"category"`
	Inventories   []Inventory    `gorm:"foreignKey:ProductId"`
	OrderDetails  []OrderDetail  `gorm:"foreignKey:ProductId"`
	ProductImages []ProductImage `gorm:"foreignKey:ProductId"`
}

type ProductImage struct {
	ID         uint      `gorm:"primarykey;autoIncrement"`
	ProductId  uint      `json:"productId"`
	ImageUrl   string    `json:"imageUrl"`
	FileName   string    `json:"fileName"`
	UploadedAt time.Time `json:"uploadedAt"`
}

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       uint   `json:"price" binding:"required"`
	Category    string `json:"category" binding:"required"`
}

type UpdateProductRequest struct {
	Name        string `json:"name" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
	Price       uint   `json:"price" binding:"omitempty"`
	Category    string `json:"category" binding:"omitempty"`
}

type ProductResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ProductDataResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       uint      `json:"price"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
