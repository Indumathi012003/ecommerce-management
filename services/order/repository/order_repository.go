package repository

import (
	"context"
	"order-management/model"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) (*model.Order, error)
	GetByID(ctx context.Context, id int) (*model.Order, error)
	ListAll(ctx context.Context) ([]*model.Order, error)
	UpdateStatus(ctx context.Context, id int, status model.OrderStatus) (*model.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *model.Order) (*model.Order, error) {
	err := r.db.WithContext(ctx).Create(order).Error
	return order, err
}

func (r *orderRepository) GetByID(ctx context.Context, id int) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) ListAll(ctx context.Context) ([]*model.Order, error) {
	var orders []*model.Order
	err := r.db.WithContext(ctx).Order("id asc").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id int, status model.OrderStatus) (*model.Order, error) {
	var order model.Order
	if err := r.db.WithContext(ctx).First(&order, id).Error; err != nil {
		return nil, err
	}
	order.Status = status
	err := r.db.WithContext(ctx).Save(&order).Error
	return &order, err
}