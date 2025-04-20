package entity

import (
	"time"
)

type Method string

const (
	Transfer Method = "transfer"
	Voucher  Method = "voucher"
)

type Payment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	OrderId   uint      `gorm:"not null" json:"orderId"`
	Method    Method    `gorm:"not null;type:enum('transfer', 'voucher')" json:"method"`
	Amount    uint      `gorm:"not null" json:"amount"`
	Status    uint      `gorm:"not null" json:"status"`
	CreatedAt time.Time `gorm:"not null" json:"createdAt"`
}

type PaymentRequest struct {
	OrderId uint   `json:"orderId"`
	Method  Method `json:"method"`
	Amount  uint   `json:"amount"`
}

type PaymentDataResponse struct {
	ID        uint      `json:"id"`
	OrderId   uint      `json:"orderId"`
	Method    Method    `json:"method"`
	Amount    uint      `json:"amount"`
	Status    uint      `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}
