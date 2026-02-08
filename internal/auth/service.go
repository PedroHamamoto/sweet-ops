package auth

import (
	"context"
	"errors"
	"sweet-ops/internal/user"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type Service struct {
	store       *Store
	jwt         *Jwt
	userService *user.Service
}

func NewService(store *Store, jwt *Jwt, userService *user.Service) *Service {
	return &Service{
		store:       store,
		jwt:         jwt,
		userService: userService,
	}
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken uuid.UUID
}

func (s *Service) Login(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	user, err := s.userService.GetByEmail(ctx, input.Email)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)) != nil {
		return nil, ErrInvalidCredentials
	}

	// TODO: execute the following two operations in a transaction
	accessToken, err := s.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refresh, _ := uuid.NewV7()

	err = s.store.SaveRefreshToken(ctx, user.ID, refresh.String(), time.Now().Add(30*24*time.Hour))
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refresh,
	}, nil
}
