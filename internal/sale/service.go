package sale

import (
	"context"
	"fmt"
	"sweet-ops/internal/types"
	"sweet-ops/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

type CreateSaleItemInput struct {
	ProductID uuid.UUID
	Quantity  int
	IsGift    bool
}

type CreateSaleInput struct {
	Source          Source
	PaymentMethod   PaymentMethod
	SelfConsumption bool
	Items           []CreateSaleItemInput
}

func (s *Service) Create(ctx context.Context, input *CreateSaleInput) (*Sale, error) {
	if input.SelfConsumption {
		input.Source = SourceSelfConsumption
		input.PaymentMethod = PaymentSelfConsumption
	} else {
		if !input.Source.IsValid() {
			return nil, ErrInvalidSource
		}
		if !input.PaymentMethod.IsValid() {
			return nil, ErrInvalidPaymentMethod
		}
	}
	if len(input.Items) == 0 {
		return nil, ErrEmptyItems
	}

	for _, item := range input.Items {
		if item.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
		if input.SelfConsumption && item.IsGift {
			return nil, ErrGiftInSelfConsumption
		}
	}

	saleID := utils.NewUUID()
	now := time.Now()

	sale := &Sale{
		ID:              saleID,
		Source:          input.Source,
		PaymentMethod:   input.PaymentMethod,
		SelfConsumption: input.SelfConsumption,
		CreatedAt:       now,
	}

	err := utils.ExecuteTx(ctx, s.store, func(tx pgx.Tx) error {
		var items []*SaleItem
		var total float64

		for _, itemInput := range input.Items {
			stock, sellingPrice, version, err := s.store.GetProductStockAndPrice(ctx, tx, itemInput.ProductID)
			if err != nil {
				return fmt.Errorf("product %s: %w", itemInput.ProductID, err)
			}

			if stock < itemInput.Quantity {
				return fmt.Errorf("product %s: %w", itemInput.ProductID, ErrInsufficientStock)
			}

			unitPrice := sellingPrice
			if itemInput.IsGift {
				unitPrice = 0
			}

			item := &SaleItem{
				ID:        utils.NewUUID(),
				SaleID:    saleID,
				ProductID: itemInput.ProductID,
				Quantity:  itemInput.Quantity,
				UnitPrice: unitPrice,
				IsGift:    itemInput.IsGift,
			}
			items = append(items, item)

			if !itemInput.IsGift && !input.SelfConsumption {
				total += unitPrice * float64(itemInput.Quantity)
			}

			if err := s.store.DecrementStock(ctx, tx, itemInput.ProductID, itemInput.Quantity, version); err != nil {
				return fmt.Errorf("product %s: %w", itemInput.ProductID, err)
			}
		}

		sale.Total = total
		sale.Items = items

		if err := s.store.SaveSale(ctx, tx, sale); err != nil {
			return err
		}

		for _, item := range items {
			if err := s.store.SaveSaleItem(ctx, tx, item); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return sale, nil
}

func (s *Service) GetAll(ctx context.Context, page, pageSize int) (types.Pageable[*Sale], error) {
	offset := (page - 1) * pageSize
	sales, totalItems, err := s.store.FindAll(ctx, pageSize, offset)
	if err != nil {
		return types.Pageable[*Sale]{}, err
	}

	for _, sale := range sales {
		items, err := s.store.FindItemsBySaleID(ctx, sale.ID)
		if err != nil {
			return types.Pageable[*Sale]{}, err
		}
		sale.Items = items
	}

	return types.NewPageable(sales, page, pageSize, totalItems), nil
}
