package repository

import (
	"context"
	"order-management/model"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, p *model.Product) (*model.Product, error)
	GetByID(ctx context.Context, id int) (*model.Product, error)
	ListAll(ctx context.Context) ([]*model.Product, error)
	Update(ctx context.Context, id int, p *model.Product) (*model.Product, error)
	Delete(ctx context.Context, id int) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, p *model.Product) (*model.Product, error) {
	err := r.db.WithContext(ctx).Create(p).Error
	return p, err
}

func (r *productRepository) GetByID(ctx context.Context, id int) (*model.Product, error) {
	var p model.Product
	err := r.db.WithContext(ctx).First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) ListAll(ctx context.Context) ([]*model.Product, error) {
	var products []*model.Product
	err := r.db.WithContext(ctx).Order("id asc").Find(&products).Error
	return products, err
}

func (r *productRepository) Update(ctx context.Context, id int, p *model.Product) (*model.Product, error) {
	var existing model.Product
	if err := r.db.WithContext(ctx).First(&existing, id).Error; err != nil {
		return nil, err
	}
	existing.Name = p.Name
	existing.Price = p.Price
	existing.Stock = p.Stock
	err := r.db.WithContext(ctx).Save(&existing).Error
	return &existing, err
}

func (r *productRepository) Delete(ctx context.Context, id int) error {
	res := r.db.WithContext(ctx).Delete(&model.Product{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}