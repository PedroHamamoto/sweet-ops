package category

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sweet-ops/internal/types"
	"sweet-ops/internal/ui"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type CategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) Create(w http.ResponseWriter, req *http.Request) {
	var request CreateCategoryRequest

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		// TODO return a proper error response
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	category, err := h.service.Create(req.Context(), request.Name)
	if err != nil {
		// TODO return a proper error response
		http.Error(w, "failed to create category", http.StatusInternalServerError)
		return
	}

	response := CategoryResponse{
		ID:   category.ID.String(),
		Name: category.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RenderCategories(w http.ResponseWriter, req *http.Request) {
	ui.Render(w, req, "categories", nil)
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
		// TODO return a proper error response
		http.Error(w, "failed to get categories", http.StatusInternalServerError)
		return
	}

	var items []CategoryResponse
	for _, category := range result.Data {
		items = append(items, CategoryResponse{
			ID:   category.ID.String(),
			Name: category.Name,
		})
	}

	response := types.NewPageable(items, result.Page, result.PageSize, result.TotalItems)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
