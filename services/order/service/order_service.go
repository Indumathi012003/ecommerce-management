package service

import (
	"context"
	"gorm.io/gorm"
	"errors"
	"fmt"
	"time"

	"order-management/model"
	"order-management/services/order/client"
	"order-management/services/order/repository"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req model.CreateOrderRequest) (*model.Order, error)
	GetOrder(ctx context.Context, id int) (*model.Order, error)
	ListOrders(ctx context.Context) ([]*model.Order, error)
	UpdateOrderStatus(ctx context.Context, id int, status model.OrderStatus) (*model.Order, error)
}

type orderService struct {
	repo          repository.OrderRepository
	productClient client.ProductClient
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{
		repo:          repo,
		productClient: client.NewProductClient(),
	}
}

func NewOrderServiceWithClient(repo repository.OrderRepository, pc client.ProductClient) OrderService {
	return &orderService{
		repo:          repo,
		productClient: pc,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, req model.CreateOrderRequest) (*model.Order, error) {
	// 1. Fetch live product price via resilient client
	price, err := s.productClient.GetProductPrice(ctx, int64(req.ProductID))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}

	// 2. Compute total amount
	totalAmount := price * float64(req.Quantity)
	order := &model.Order{
		CustomerID:  req.CustomerID,
		ProductID:   req.ProductID,
		Quantity:    req.Quantity,
		TotalAmount: totalAmount,
		Status:      model.StatusConfirmed,
		CreatedAt:   time.Now(),
	}

	// 3. Save order to database
	createdOrder, err := s.repo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return createdOrder, nil
}

func (s *orderService) GetOrder(ctx context.Context, id int) (*model.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *orderService) ListOrders(ctx context.Context) ([]*model.Order, error) {
	return s.repo.ListAll(ctx)
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, id int, status model.OrderStatus) (*model.Order, error) {
	if status != model.StatusCreated && status != model.StatusConfirmed && status != model.StatusCancelled {
	// if status != model.StatusCreated && status != model.StatusConfirmed && status != model.StatusCancelled && status != model.StatusRefunded {
		return nil, errors.New("invalid order status")
	}
	return s.repo.UpdateStatus(ctx, id, status)
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
