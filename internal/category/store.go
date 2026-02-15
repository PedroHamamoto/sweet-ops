package category

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrCategoryNotFound = errors.New("category not found")

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
	err := s.db.QueryRow(ctx, statement, id, name).
		Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *Store) FindAll(ctx context.Context, page, pageSize int) ([]*Category, int, error) {
	var totalItems int
	countStmt := "SELECT COUNT(*) FROM categories"
	if err := s.db.QueryRow(ctx, countStmt).Scan(&totalItems); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	statement := "SELECT id, name, created_at, updated_at FROM categories ORDER BY id LIMIT $1 OFFSET $2"
	rows, err := s.db.Query(ctx, statement, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		c := &Category{}
		err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return categories, totalItems, nil
}

func (s *Store) FindByID(ctx context.Context, id uuid.UUID) (*Category, error) {
	statement := "SELECT id, name, created_at, updated_at FROM categories WHERE id = $1"
	c := &Category{}
	err := s.db.QueryRow(ctx, statement, id).Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return c, nil
}
