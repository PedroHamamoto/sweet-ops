package router

import (
	"net/http"
	"sweet-ops/internal/auth"
	"sweet-ops/internal/user"

	"github.com/go-chi/chi/v5"
)

func NewRouter(userHandler *user.Handler, authHandler *auth.Handler) http.Handler {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Post("/users", userHandler.Create)
		r.Post("/login", authHandler.Login)
	})

	return r
}
