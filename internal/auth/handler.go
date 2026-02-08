package auth

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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) Login(w http.ResponseWriter, req *http.Request) {
	var request LoginRequest

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		log.Printf("failed to decode request body: %v", err)
		// TODO: return a propper error response
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	loginInput := &LoginInput{
		Email:    request.Email,
		Password: request.Password,
	}

	output, err := h.service.Login(req.Context(), loginInput)
	if err != nil {
		log.Printf("failed to login: %v", err)
		// TODO: return a propper error response
		// TODO: Improve error handling
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	response := LoginResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken.String(),
	}

	// TODO: Use a middleware to set the content type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
