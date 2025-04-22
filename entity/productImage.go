package entity

import "time"

type ProductImage struct {
	ID         uint      `gorm:"primarykey;autoIncrement"`
	ProductId  uint      `json:"productId"`
	ImageUrl   string    `gorm:"type:text" json:"imageUrl"`
	FileName   string    `json:"fileName"`
	UploadedAt time.Time `json:"uploadedAt"`
}
