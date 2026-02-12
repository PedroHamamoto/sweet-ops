package category

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, name string) (*Category, error) {
	id, _ := uuid.NewV7()
	statement := "INSERT INTO categories (id, name) VALUES ($1, $2) RETURNING id, name, created_at, updated_at"
	c := &Category{}
	err := s.db.QueryRow(context.Background(), statement, id, name).
		Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return c, nil
}
