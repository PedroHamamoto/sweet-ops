package product

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

func (s *Store) Create(ctx context.Context, product *Product) (*Product, error) {
	statement := `
		INSERT INTO products (id, category_id, flavor, production_price, selling_price, markup_margin, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, current_timestamp, current_timestamp)
		RETURNING id, created_at, updated_at
		`

	id, _ := uuid.NewV7()

	err := s.db.QueryRow(ctx, statement, id, product.Category.ID, product.Flavor, product.ProductionPrice, product.SellingPrice, product.MarkupMargin).
		Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return product, nil
}
