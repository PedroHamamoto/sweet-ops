package router

import (
	"net/http"
	"sweet-ops/internal/auth"
	"sweet-ops/internal/ui"
	"sweet-ops/internal/user"

	"github.com/go-chi/chi/v5"
)

func NewRouter(userHandler *user.Handler, authHandler *auth.Handler, authMiddleware *auth.Middleware) http.Handler {
	r := chi.NewRouter()

	// Pages
	r.Group(func(r chi.Router) {
		r.Use(ui.Layout)

		r.Get("/login", authHandler.RenderLogin)
		r.Get("/home", Home)
	})

	// API
	r.Route("/api", func(r chi.Router) {
		r.Post("/users", userHandler.Create)
		r.Post("/login", authHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("This is a protected route"))
			})
		})
	})

	return r
}

// TODO: Move it to a better place when we have more pages and handlers
func Home(w http.ResponseWriter, r *http.Request) {
	ui.Render(w, r, "home", nil)
}
