package service

import (
	"context"
	"gorm.io/gorm"
	"errors"
	"time"

	"order-management/model"
	"order-management/services/product/repository"
)

type ProductService interface {
	AddProduct(ctx context.Context, req model.CreateProductRequest) (*model.Product, error)
	GetProduct(ctx context.Context, id int) (*model.Product, error)
	ListProducts(ctx context.Context) ([]*model.Product, error)
	UpdateProduct(ctx context.Context, id int, req model.UpdateProductRequest) (*model.Product, error)
	DeleteProduct(ctx context.Context, id int) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) AddProduct(ctx context.Context, req model.CreateProductRequest) (*model.Product, error) {
	product := &model.Product{
		Name:      req.Name,
		Price:     req.Price,
		Stock:     req.Stock,
		CreatedAt: time.Now(),
	}
	return s.repo.Create(ctx, product)
}

func (s *productService) GetProduct(ctx context.Context, id int) (*model.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *productService) ListProducts(ctx context.Context) ([]*model.Product, error) {
	return s.repo.ListAll(ctx)
}

func (s *productService) UpdateProduct(ctx context.Context, id int, req model.UpdateProductRequest) (*model.Product, error) {
	product := &model.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}
	return s.repo.Update(ctx, id, product)
}

func (s *productService) DeleteProduct(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
