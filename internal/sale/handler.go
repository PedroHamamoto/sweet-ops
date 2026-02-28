package sale

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sweet-ops/internal/types"
	"sweet-ops/internal/ui"
	"time"

	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type CreateSaleItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	IsGift    bool   `json:"is_gift"`
}

type CreateSaleRequest struct {
	Source          string                  `json:"source"`
	PaymentMethod   string                  `json:"payment_method"`
	SelfConsumption bool                    `json:"self_consumption"`
	Items           []CreateSaleItemRequest `json:"items"`
}

type SaleItemResponse struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	IsGift      bool    `json:"is_gift"`
}

type SaleResponse struct {
	ID              string             `json:"id"`
	Source          string             `json:"source"`
	PaymentMethod   string             `json:"payment_method"`
	SelfConsumption bool               `json:"self_consumption"`
	Total           float64            `json:"total"`
	Items           []SaleItemResponse `json:"items"`
	CreatedAt       time.Time          `json:"created_at"`
}

func (h *Handler) Create(w http.ResponseWriter, req *http.Request) {
	var request CreateSaleRequest

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var items []CreateSaleItemInput
	for _, item := range request.Items {
		productID, err := uuid.Parse(item.ProductID)
		if err != nil {
			http.Error(w, "invalid product id: "+item.ProductID, http.StatusBadRequest)
			return
		}
		items = append(items, CreateSaleItemInput{
			ProductID: productID,
			Quantity:  item.Quantity,
			IsGift:    item.IsGift,
		})
	}

	input := &CreateSaleInput{
		Source:          Source(request.Source),
		PaymentMethod:   PaymentMethod(request.PaymentMethod),
		SelfConsumption: request.SelfConsumption,
		Items:           items,
	}

	sale, err := h.service.Create(req.Context(), input)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrInvalidSource) || errors.Is(err, ErrInvalidPaymentMethod) ||
			errors.Is(err, ErrEmptyItems) || errors.Is(err, ErrInvalidQuantity) ||
			errors.Is(err, ErrGiftInSelfConsumption) || errors.Is(err, ErrInsufficientStock) {
			status = http.StatusBadRequest
		}
		http.Error(w, err.Error(), status)
		return
	}

	response := toSaleResponse(sale)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetAll(w http.ResponseWriter, req *http.Request) {
	page, _ := strconv.Atoi(req.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(req.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	result, err := h.service.GetAll(req.Context(), page, pageSize)
	if err != nil {
		http.Error(w, "failed to get sales", http.StatusInternalServerError)
		return
	}

	var items []SaleResponse
	for _, sale := range result.Data {
		items = append(items, toSaleResponse(sale))
	}

	response := types.NewPageable(items, result.Page, result.PageSize, result.TotalItems)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RenderSales(w http.ResponseWriter, req *http.Request) {
	ui.Render(w, req, "sales", nil)
}

func toSaleResponse(sale *Sale) SaleResponse {
	var itemResponses []SaleItemResponse
	for _, item := range sale.Items {
		itemResponses = append(itemResponses, SaleItemResponse{
			ID:          item.ID.String(),
			ProductID:   item.ProductID.String(),
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			IsGift:      item.IsGift,
		})
	}

	return SaleResponse{
		ID:              sale.ID.String(),
		Source:          string(sale.Source),
		PaymentMethod:   string(sale.PaymentMethod),
		SelfConsumption: sale.SelfConsumption,
		Total:           sale.Total,
		Items:           itemResponses,
		CreatedAt:       sale.CreatedAt,
	}
}
