package product

import (
	"context"
	"sweet-ops/internal/category"

	"github.com/google/uuid"
)

type Service struct {
	store           *Store
	categoryService *category.Service
}

func NewService(store *Store, categoryService *category.Service) *Service {
	return &Service{store: store, categoryService: categoryService}
}

type CreateProductInput struct {
	CategoryID      uuid.UUID
	Flavor          string
	ProductionPrice float64
	SellingPrice    float64
}

func (s *Service) Create(ctx context.Context, input *CreateProductInput) (*Product, error) {
	category, err := s.categoryService.GetByID(ctx, input.CategoryID)
	if err != nil {
		return nil, err
	}

	product := NewProduct(category, input.Flavor, input.ProductionPrice, input.SellingPrice)

	return s.store.Create(ctx, product)
}
