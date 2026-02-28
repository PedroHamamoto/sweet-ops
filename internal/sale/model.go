package sale

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSaleNotFound          = errors.New("sale not found")
	ErrEmptyItems            = errors.New("a sale must contain at least one item")
	ErrInvalidQuantity       = errors.New("quantity must be greater than 0")
	ErrInvalidSource         = errors.New("invalid source")
	ErrInvalidPaymentMethod  = errors.New("invalid payment method")
	ErrGiftInSelfConsumption = errors.New("gift items are not allowed in self-consumption sales")
	ErrInsufficientStock     = errors.New("insufficient stock for product")
	ErrVersionMismatch       = errors.New("optimistic lock error: product was updated by another transaction")
)

type Source string

const (
	SourceBalcao Source = "Balcão"
	SourceGoomer Source = "Goomer"
	Source99Food Source = "99 Food"
)

func (s Source) IsValid() bool {
	switch s {
	case SourceBalcao, SourceGoomer, Source99Food:
		return true
	}
	return false
}

type PaymentMethod string

const (
	PaymentCreditCard PaymentMethod = "Cartão de Crédito"
	PaymentDebitCard  PaymentMethod = "Cartão de Débito"
	PaymentPIX        PaymentMethod = "PIX"
	PaymentCash       PaymentMethod = "Dinheiro"
)

func (p PaymentMethod) IsValid() bool {
	switch p {
	case PaymentCreditCard, PaymentDebitCard, PaymentPIX, PaymentCash:
		return true
	}
	return false
}

type Sale struct {
	ID              uuid.UUID
	Items           []*SaleItem
	Source          Source
	PaymentMethod   PaymentMethod
	SelfConsumption bool
	Total           float64
	CreatedAt       time.Time
}

type SaleItem struct {
	ID          uuid.UUID
	SaleID      uuid.UUID
	ProductID   uuid.UUID
	Quantity    int
	UnitPrice   float64
	IsGift      bool
	ProductName string
}
