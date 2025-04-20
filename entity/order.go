package entity

import "time"

type Order struct {
	ID              uint          `gorm:"primaryKey;autoIncrement"`
	UserId          uint          `json:"userId"`
	Date            time.Time     `gorm:"not null" json:"date"`
	ShippingAddress string        `gorm:"not null" json:"shippingAddress"`
	Total           uint          `gorm:"not null" json:"total"`
	Status          uint          `gorm:"not null" json:"status"`
	OrderDetails    []OrderDetail `gorm:"foreignKey:OrderId"`
	Payment         Payment       `gorm:"foreignKey:OrderId"`
}

type OrderDetail struct {
	OrderId   uint `gorm:"primaryKey"`
	ProductId uint `gorm:"primaryKey"`
	Qty       uint `json:"qty"`
	Subtotal  uint `json:"subtotal"`
}

type CreateOrderRequest struct {
	ShippingAddress string               `json:"shippingAddress"`
	OrderDetails    []OrderDetailRequest `json:"orderDetails" binding:"required,dive"`
}

type OrderDetailRequest struct {
	ProductId uint `json:"productId" binding:"required"`
	Qty       uint `json:"qty" binding:"required"`
}

type OrderDataResponse struct {
	ID                  uint                  `json:"id"`
	UserId              uint                  `json:"userId"`
	Date                time.Time             `json:"date"`
	ShippingAddress     string                `json:"shippingAddress"`
	Total               uint                  `json:"total"`
	Status              uint                  `json:"status"`
	OrderDetailResponse []OrderDetailResponse `json:"orderDetails"`
}

type OrderDetailResponse struct {
	ProductId uint `json:"productId"`
	Qty       uint `json:"qty"`
	Subtotal  uint `json:"subtotal"`
}
