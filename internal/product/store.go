package product

import (
	"context"
	"errors"
	"sweet-ops/internal/category"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) DB() *pgxpool.Pool {
	return s.db
}

func (s *Store) Create(ctx context.Context, product *Product) (*Product, error) {
	statement := `
		INSERT INTO products (id, category_id, flavor, production_price, selling_price, markup_margin, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, current_timestamp, current_timestamp)
		RETURNING id, created_at, updated_at
		`

	err := s.db.QueryRow(ctx, statement, product.ID, product.Category.ID, product.Flavor, product.ProductionPrice, product.SellingPrice, product.MarkupMargin).
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

func (s *Store) GetVersion(ctx context.Context, tx pgx.Tx, productID uuid.UUID) (int, error) {
	var version int
	err := tx.QueryRow(ctx, "SELECT version FROM products WHERE id = $1", productID).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrProductNotFound
		}
		return 0, err
	}
	return version, nil
}

func (s *Store) SaveProduction(ctx context.Context, tx pgx.Tx, production *Production) error {
	_, err := tx.Exec(ctx, "INSERT INTO productions (id, product_id, quantity, created_at) VALUES ($1, $2, $3, $4)",
		production.ID, production.ProductID, production.Quantity, production.CreatedAt)
	return err
}

func (s *Store) IncrementStock(ctx context.Context, tx pgx.Tx, productID uuid.UUID, quantity int, currentVersion int) error {
	tag, err := tx.Exec(ctx, "UPDATE products SET stock_quantity = stock_quantity + $2, version = version + 1 WHERE id = $1 AND version = $3",
		productID, quantity, currentVersion)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("optimistic lock error: product was updated by another transaction")
	}

	return nil
}
