package product

import (
	"sweet-ops/internal/category"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID              uuid.UUID
	Category        *category.Category
	Flavor          string
	ProductionPrice float64
	SellingPrice    float64
	MarkupMargin    float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewProduct(category *category.Category, flavor string, productionPrice, sellingPrice float64) *Product {
	markupMargin := ((sellingPrice - productionPrice) / productionPrice) * 100
	return &Product{
		Category:        category,
		Flavor:          flavor,
		ProductionPrice: productionPrice,
		SellingPrice:    sellingPrice,
		MarkupMargin:    markupMargin,
	}
}
