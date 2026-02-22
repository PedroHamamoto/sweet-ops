package product

import (
	"context"
	"sweet-ops/internal/category"
	"sweet-ops/internal/types"
	"sweet-ops/internal/utils"

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

	id := utils.NewUUID()
	product := NewProduct(id, category, input.Flavor, input.ProductionPrice, input.SellingPrice)

	return s.store.Create(ctx, product)
}

func (s *Service) GetAll(ctx context.Context, page, pageSize int) (types.Pageable[*Product], error) {
	products, totalItems, err := s.store.FindAll(ctx, page, pageSize)
	if err != nil {
		return types.Pageable[*Product]{}, err
	}
	return types.NewPageable(products, page, pageSize, totalItems), nil
}

func (s *Service) RegisterProduction(ctx context.Context, productID uuid.UUID, quantity int) error {
	id := utils.NewUUID()
	return s.store.RegisterProduction(ctx, productID, id, quantity)
}
