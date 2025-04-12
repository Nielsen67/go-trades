package entity

import "time"

type Order struct {
	ID              uint          `gorm:"primaryKey;autoIncrement"`
	UserId          uint          `json:"userId"`
	Date            time.Time     `gorm:"not null" json:"date"`
	ShippingAddress string        `gorm:"not null" json:"shippingAddress"`
	ShippingCity    string        `gorm:"not null" json:"shippingCity"`
	Total           uint          `gorm:"not null" json:"total"`
	Status          string        `gorm:"not null" json:"status"`
	OrderDetails    []OrderDetail `gorm:"foreignKey:OrderId"`
	Payment         Payment       `gorm:"foreignKey:OrderId"`
}

type OrderDetail struct {
	OrderId     uint `gorm:"primaryKey"`
	ProductId   uint `gorm:"primaryKey"`
	Qty         uint `json:"qty"`
	PriceDetail uint `json:"priceDetail"`
}

type CreateOrderRequest struct {
	ShippingAddress string               `json:"shippingAddress"`
	ShippingCity    string               `json:"shippingCity"`
	OrderDetails    []OrderDetailRequest `json:"orderDetails" binding:"required,dive"`
}

type OrderDetailRequest struct {
	ProductId uint `json:"productId" binding:"required"`
	Qty       uint `json:"qty" binding:"required"`
}

type OrderResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type OrderDataResponse struct {
	ID                  uint                  `json:"id"`
	UserId              uint                  `json:"userId"`
	Date                time.Time             `json:"date"`
	ShippingAddress     string                `json:"shippingAddress"`
	ShippingCity        string                `json:"shippingCity"`
	Total               uint                  `json:"total"`
	Status              string                `json:"status"`
	OrderDetailResponse []OrderDetailResponse `json:"orderDetails"`
}

type OrderDetailResponse struct {
	ProductId   uint `json:"productId"`
	Qty         uint `json:"qty"`
	PriceDetail uint `json:"priceDetail"`
}
