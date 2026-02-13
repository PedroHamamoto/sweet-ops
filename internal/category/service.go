package category

import (
	"context"
	"sweet-ops/internal/types"
)

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(ctx context.Context, name string) (*Category, error) {
	return s.store.Create(ctx, name)
}

func (s *Service) GetAll(ctx context.Context, page, pageSize int) (types.Pageable[*Category], error) {
	categories, totalItems, err := s.store.FindAll(ctx, page, pageSize)
	if err != nil {
		return types.Pageable[*Category]{}, err
	}
	return types.NewPageable(categories, page, pageSize, totalItems), nil
}
