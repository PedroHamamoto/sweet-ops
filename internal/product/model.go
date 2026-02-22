package product

import (
	"errors"
	"sweet-ops/internal/category"
	"time"

	"github.com/google/uuid"
)

var ErrProductNotFound = errors.New("product not found")

type Product struct {
	ID              uuid.UUID
	Category        *category.Category
	Flavor          string
	ProductionPrice float64
	SellingPrice    float64
	MarkupMargin    float64
	StockQuantity   int
	Version         int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Production struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	Quantity  int
	CreatedAt time.Time
}

func NewProduct(id uuid.UUID, category *category.Category, flavor string, productionPrice, sellingPrice float64) *Product {
	markupMargin := ((sellingPrice - productionPrice) / productionPrice) * 100
	return &Product{
		ID:              id,
		Category:        category,
		Flavor:          flavor,
		ProductionPrice: productionPrice,
		SellingPrice:    sellingPrice,
		MarkupMargin:    markupMargin,
	}
}
