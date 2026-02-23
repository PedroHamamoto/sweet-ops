package product

import (
	"context"
	"sweet-ops/internal/category"
	"sweet-ops/internal/types"
	"sweet-ops/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	offset := (page - 1) * pageSize
	products, totalItems, err := s.store.FindAll(ctx, pageSize, offset)
	if err != nil {
		return types.Pageable[*Product]{}, err
	}
	return types.NewPageable(products, page, pageSize, totalItems), nil
}

func (s *Service) GetAllProductions(ctx context.Context, page, pageSize int) (types.Pageable[*Production], error) {
	offset := (page - 1) * pageSize
	productions, totalItems, err := s.store.FindAllProductions(ctx, pageSize, offset)
	if err != nil {
		return types.Pageable[*Production]{}, err
	}
	return types.NewPageable(productions, page, pageSize, totalItems), nil
}

func (s *Service) RegisterProduction(ctx context.Context, productID uuid.UUID, quantity int) error {
	id := utils.NewUUID()

	return utils.ExecuteTx(ctx, s.store, func(tx pgx.Tx) error {
		version, err := s.store.GetVersion(ctx, tx, productID)
		if err != nil {
			return err
		}

		production := &Production{
			ID:        id,
			ProductID: productID,
			Quantity:  quantity,
			CreatedAt: time.Now(),
		}

		if err := s.store.SaveProduction(ctx, tx, production); err != nil {
			return err
		}

		return s.store.IncrementStock(ctx, tx, productID, quantity, version)
	})
}
