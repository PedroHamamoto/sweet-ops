package sale

import (
	"context"
	"errors"

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

func (s *Store) SaveSale(ctx context.Context, tx pgx.Tx, sale *Sale) error {
	_, err := tx.Exec(ctx,
		"INSERT INTO sales (id, source, payment_method, self_consumption, total, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		sale.ID, sale.Source, sale.PaymentMethod, sale.SelfConsumption, sale.Total, sale.CreatedAt,
	)
	return err
}

func (s *Store) SaveSaleItem(ctx context.Context, tx pgx.Tx, item *SaleItem) error {
	_, err := tx.Exec(ctx,
		"INSERT INTO sale_items (id, sale_id, product_id, quantity, unit_price, is_gift) VALUES ($1, $2, $3, $4, $5, $6)",
		item.ID, item.SaleID, item.ProductID, item.Quantity, item.UnitPrice, item.IsGift,
	)
	return err
}

func (s *Store) GetProductStockAndPrice(ctx context.Context, tx pgx.Tx, productID uuid.UUID) (stock int, sellingPrice float64, version int, err error) {
	err = tx.QueryRow(ctx,
		"SELECT stock_quantity, selling_price, version FROM products WHERE id = $1",
		productID,
	).Scan(&stock, &sellingPrice, &version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, 0, errors.New("product not found")
		}
	}
	return
}

func (s *Store) DecrementStock(ctx context.Context, tx pgx.Tx, productID uuid.UUID, quantity int, currentVersion int) error {
	tag, err := tx.Exec(ctx,
		"UPDATE products SET stock_quantity = stock_quantity - $2, version = version + 1 WHERE id = $1 AND version = $3 AND stock_quantity >= $2",
		productID, quantity, currentVersion,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		var version int
		err := tx.QueryRow(ctx, "SELECT version FROM products WHERE id = $1", productID).Scan(&version)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.New("product not found")
			}
			return err
		}
		if version != currentVersion {
			return ErrVersionMismatch
		}
		return ErrInsufficientStock
	}
	return nil
}

func (s *Store) FindAll(ctx context.Context, limit, offset int) ([]*Sale, int, error) {
	var totalItems int
	if err := s.db.QueryRow(ctx, "SELECT COUNT(*) FROM sales").Scan(&totalItems); err != nil {
		return nil, 0, err
	}

	statement := `
		SELECT id, source, payment_method, self_consumption, total, created_at
		FROM sales
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := s.db.Query(ctx, statement, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sales []*Sale
	for rows.Next() {
		sale := &Sale{}
		if err := rows.Scan(&sale.ID, &sale.Source, &sale.PaymentMethod, &sale.SelfConsumption, &sale.Total, &sale.CreatedAt); err != nil {
			return nil, 0, err
		}
		sales = append(sales, sale)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return sales, totalItems, nil
}

func (s *Store) FindItemsBySaleID(ctx context.Context, saleID uuid.UUID) ([]*SaleItem, error) {
	statement := `
		SELECT si.id, si.sale_id, si.product_id, si.quantity, si.unit_price, si.is_gift,
		       c.name || ' ' || p.flavor AS product_name
		FROM sale_items si
		JOIN products p ON si.product_id = p.id
		JOIN categories c ON p.category_id = c.id
		WHERE si.sale_id = $1
		ORDER BY si.is_gift, product_name
	`
	rows, err := s.db.Query(ctx, statement, saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*SaleItem
	for rows.Next() {
		item := &SaleItem{}
		if err := rows.Scan(&item.ID, &item.SaleID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.IsGift, &item.ProductName); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
