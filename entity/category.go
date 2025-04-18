package entity

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Code     string    `gorm:"unique;not null" json:"code"`
	Name     string    `gorm:"unique;not null" json:"name"`
	Products []Product `gorm:"foreignKey:CategoryId"`

}

type CategoryRequest struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type CategoryDataResponse struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
