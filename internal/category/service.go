package category

import (
	"context"
	"sweet-ops/internal/types"
	"sweet-ops/internal/utils"

	"github.com/google/uuid"
)

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(ctx context.Context, name string) (*Category, error) {
	category := &Category{
		ID:   utils.NewUUID(),
		Name: name,
	}
	return s.store.Create(ctx, category)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Category, error) {
	return s.store.FindByID(ctx, id)
}

func (s *Service) GetAll(ctx context.Context, page, pageSize int) (types.Pageable[*Category], error) {
	categories, totalItems, err := s.store.FindAll(ctx, page, pageSize)
	if err != nil {
		return types.Pageable[*Category]{}, err
	}
	return types.NewPageable(categories, page, pageSize, totalItems), nil
}
