package model

import "time"

type OrderStatus string

const (
	StatusCreated   OrderStatus = "CREATED"
	StatusConfirmed OrderStatus = "CONFIRMED"
	StatusCancelled OrderStatus = "CANCELLED"
	// StatusRefunded  OrderStatus = "REFUNDED"
)

type Order struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	CustomerID  int         `gorm:"not null;index" json:"customerId" binding:"required"`
	ProductID   int         `gorm:"not null" json:"productId" binding:"required"`
	Quantity    int         `gorm:"not null" json:"quantity" binding:"required,gt=0"`
	TotalAmount float64     `gorm:"type:numeric(10,2);not null" json:"totalAmount"`
	Status      OrderStatus `gorm:"size:50;not null;default:'CREATED'" json:"status"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

type CreateOrderRequest struct {
	CustomerID int `json:"customerId" binding:"required"`
	ProductID  int `json:"productId" binding:"required"`
	Quantity   int `json:"quantity" binding:"required,gt=0"`
}

type UpdateStatusRequest struct {
	Status OrderStatus `json:"status" binding:"required"`
}