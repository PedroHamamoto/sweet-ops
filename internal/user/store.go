package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, user *User) (*User, error) {
	statement := "INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3) RETURNING id, email, password_hash, created_at, updated_at"

	err := s.db.QueryRow(ctx, statement, user.ID, user.Email, user.PasswordHash).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) FindByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}
	statement := "SELECT id, email, password_hash FROM users WHERE email = $1"

	err := s.db.QueryRow(ctx, statement, email).Scan(&u.ID, &u.Email, &u.PasswordHash)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return u, nil

}
