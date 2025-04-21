package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	CategoryId    uint           `gorm:"not null" json:"categoryId"`
	Name          string         `gorm:"not null;unique" json:"name"`
	Description   string         `json:"description"`
	Price         uint           `gorm:"not null" json:"price"`
	Inventories   []Inventory    `gorm:"foreignKey:ProductId"`
	OrderDetails  []OrderDetail  `gorm:"foreignKey:ProductId"`
	ProductImages []ProductImage `gorm:"foreignKey:ProductId"`
}

type ProductImage struct {
	ID         uint      `gorm:"primarykey;autoIncrement"`
	ProductId  uint      `json:"productId"`
	ImageUrl   string    `gorm:"type:text" json:"imageUrl"`
	FileName   string    `json:"fileName"`
	UploadedAt time.Time `json:"uploadedAt"`
}

type CreateProductRequest struct {
	CategoryId  uint   `json:"categoryId" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       uint   `json:"price" binding:"required"`
}

type UpdateProductRequest struct {
	CategoryId  uint   `json:"categoryId" binding:"omitempty"`
	Name        string `json:"name" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
	Price       uint   `json:"price" binding:"omitempty"`
}

type ProductDataResponse struct {
	ID          uint      `json:"id"`
	CategoryId  uint      `json:"categoryId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       uint      `json:"price"`
	Stock       uint      `json:"stock"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
