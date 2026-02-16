package product

import (
	"context"
	"sweet-ops/internal/category"

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

func (s *Store) FindAll(ctx context.Context, page, pageSize int) ([]*Product, int, error) {
	var totalItems int
	countStmt := "SELECT COUNT(*) FROM products"
	if err := s.db.QueryRow(ctx, countStmt).Scan(&totalItems); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	statement := `
		SELECT p.id, p.flavor, p.production_price, p.selling_price, p.markup_margin, p.stock_quantity, p.created_at, p.updated_at,
		       c.id, c.name, c.created_at, c.updated_at
		FROM products p
		JOIN categories c ON p.category_id = c.id
		ORDER BY p.id
		LIMIT $1 OFFSET $2
	`
	rows, err := s.db.Query(ctx, statement, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		p := &Product{Category: &category.Category{}}
		err := rows.Scan(
			&p.ID, &p.Flavor, &p.ProductionPrice, &p.SellingPrice, &p.MarkupMargin, &p.StockQuantity, &p.CreatedAt, &p.UpdatedAt,
			&p.Category.ID, &p.Category.Name, &p.Category.CreatedAt, &p.Category.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}
