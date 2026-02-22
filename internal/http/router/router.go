package router

import (
	"net/http"
	"sweet-ops/internal/auth"
	"sweet-ops/internal/category"
	"sweet-ops/internal/product"
	"sweet-ops/internal/ui"
	"sweet-ops/internal/user"

	"github.com/go-chi/chi/v5"
)

func NewRouter(userHandler *user.Handler, authHandler *auth.Handler, authMiddleware *auth.Middleware, categoryHandler *category.Handler, productHandler *product.Handler) http.Handler {
	r := chi.NewRouter()

	// Static files
	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("internal/ui/static")))
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	})

	// Pages
	r.Group(func(r chi.Router) {
		r.Use(ui.Layout)

		r.Get("/login", authHandler.RenderLogin)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuthUI)
			r.Get("/home", Home)
			r.Get("/categories", categoryHandler.RenderCategories)
			r.Get("/products", productHandler.RenderProducts)
			r.Get("/productions", productHandler.RenderProductions)
		})
	})

	// API
	r.Route("/api", func(r chi.Router) {
		r.Post("/users", userHandler.Create)
		r.Post("/login", authHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Post("/categories", categoryHandler.Create)
			r.Get("/categories", categoryHandler.GetAll)

			r.Post("/products", productHandler.Create)
			r.Get("/products", productHandler.GetAll)
			r.Post("/products/{id}/productions", productHandler.RegisterProduction)
		})
	})

	return r
}

// TODO: Move it to a better place when we have more pages and handlers
func Home(w http.ResponseWriter, r *http.Request) {
	ui.Render(w, r, "home", nil)
}
