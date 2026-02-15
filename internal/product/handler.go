package product

import (
	"encoding/json"
	"errors"
	"net/http"
	"sweet-ops/internal/category"

	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type CreateProductRequest struct {
	CategoryID      string  `json:"category_id"`
	Flavor          string  `json:"flavor"`
	ProductionPrice float64 `json:"production_price"`
	SellingPrice    float64 `json:"selling_price"`
}

type ProductResponse struct {
	ID              string                    `json:"id"`
	Category        category.CategoryResponse `json:"category"`
	Flavor          string                    `json:"flavor"`
	ProductionPrice float64                   `json:"production_price"`
	SellingPrice    float64                   `json:"selling_price"`
	MarkupMargin    float64                   `json:"markup_margin"`
}

func (h *Handler) Create(w http.ResponseWriter, req *http.Request) {
	var request CreateProductRequest

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	categoryID, err := uuid.Parse(request.CategoryID)
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	input := &CreateProductInput{
		CategoryID:      categoryID,
		Flavor:          request.Flavor,
		ProductionPrice: request.ProductionPrice,
		SellingPrice:    request.SellingPrice,
	}

	product, err := h.service.Create(req.Context(), input)
	if err != nil {
		if errors.Is(err, category.ErrCategoryNotFound) {
			http.Error(w, "category not found", http.StatusBadRequest)
			return
		}
		http.Error(w, "failed to create product", http.StatusInternalServerError)
		return
	}

	response := ProductResponse{
		ID: product.ID.String(),
		Category: category.CategoryResponse{
			ID:   product.Category.ID.String(),
			Name: product.Category.Name,
		},
		Flavor:          product.Flavor,
		ProductionPrice: product.ProductionPrice,
		SellingPrice:    product.SellingPrice,
		MarkupMargin:    product.MarkupMargin,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
