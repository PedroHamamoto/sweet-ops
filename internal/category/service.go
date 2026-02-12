package category

import "context"

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(ctx context.Context, name string) (*Category, error) {
	return s.store.Create(ctx, name)
}
