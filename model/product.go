package model

import "time"

type Product struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name" binding:"required"`
	Price     float64   `gorm:"type:numeric(10,2);not null" json:"price" binding:"required,gte=0"`
	Stock     int       `gorm:"not null;default:0" json:"stock" binding:"gte=0"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateProductRequest struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gte=0"`
	Stock int     `json:"stock" binding:"gte=0"`
}

type UpdateProductRequest struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gte=0"`
	Stock int     `json:"stock" binding:"gte=0"`
}