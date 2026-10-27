package mock

import (
	"context"
	"testing"

	"order-management/model"
	"order-management/services/order/service"

	"github.com/stretchr/testify/assert"
)

// Minimal mock for OrderRepository
type mockRepo struct{}

func (m *mockRepo) Create(ctx context.Context, o *model.Order) (*model.Order, error) {
	o.ID = 1
	return o, nil
}
func (m *mockRepo) GetByID(ctx context.Context, id int) (*model.Order, error) {
    return &model.Order{ID: uint(id)}, nil
}
func (m *mockRepo) ListAll(ctx context.Context) ([]*model.Order, error) {
	return []*model.Order{}, nil
}
func (m *mockRepo) UpdateStatus(ctx context.Context, id int, s model.OrderStatus) (*model.Order, error) {
    return &model.Order{ID: uint(id), Status: s}, nil
}

// Minimal mock for ProductClient
type mockClient struct{}

func (m *mockClient) GetProductPrice(ctx context.Context, id int64) (float64, error) {
	return 100.0, nil
}

func TestOrderService_CreateOrder(t *testing.T) {
	svc := service.NewOrderServiceWithClient(&mockRepo{}, &mockClient{})

	req := model.CreateOrderRequest{
		CustomerID: 101,
		ProductID:  1,
		Quantity:   2,
	}

	order, err := svc.CreateOrder(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, 200.0, order.TotalAmount)
}