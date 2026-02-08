package user

import (
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (h *Handler) Create(w http.ResponseWriter, req *http.Request) {
	var request CreateUserRequest

	// TODO validate request body

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		// TODO return a proper error response
		log.Printf("failed to decode request body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	CreateUserInput := CreateUserInput{
		Email:    request.Email,
		Password: request.Password,
	}

	user, err := h.service.CreateUser(req.Context(), CreateUserInput)
	if err != nil {
		// TODO return a proper error response
		log.Printf("failed to create user: %v", err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	response := UserResponse{
		ID:    user.ID.String(),
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
