package user

import (
	"context"
	"sweet-ops/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

type CreateUserInput struct {
	Email    string
	Password string
}

func (s *Service) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	id := utils.NewUUID()

	user := &User{
		ID:           id,
		Email:        input.Email,
		PasswordHash: string(hash),
	}

	createdUser, err := s.store.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*User, error) {
	return s.store.FindByEmail(ctx, email)
}
