package category

import (
	"encoding/json"
	"net/http"
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
