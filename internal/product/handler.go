package product

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sweet-ops/internal/category"
	"sweet-ops/internal/types"
	"sweet-ops/internal/ui"

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

func (h *Handler) RenderProducts(w http.ResponseWriter, req *http.Request) {
	ui.Render(w, req, "products", nil)
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
		http.Error(w, "failed to get products", http.StatusInternalServerError)
		return
	}

	var items []ProductResponse
	for _, product := range result.Data {
		items = append(items, ProductResponse{
			ID: product.ID.String(),
			Category: category.CategoryResponse{
				ID:   product.Category.ID.String(),
				Name: product.Category.Name,
			},
			Flavor:          product.Flavor,
			ProductionPrice: product.ProductionPrice,
			SellingPrice:    product.SellingPrice,
			MarkupMargin:    product.MarkupMargin,
		})
	}

	response := types.NewPageable(items, result.Page, result.PageSize, result.TotalItems)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
