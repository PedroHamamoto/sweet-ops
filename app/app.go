package app

import (
	"net/http"
	"os"
	"sweet-ops/internal/auth"
	"sweet-ops/internal/category"
	"sweet-ops/internal/http/router"
	"sweet-ops/internal/product"
	"sweet-ops/internal/sale"
	"sweet-ops/internal/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Router http.Handler
}

func New(db *pgxpool.Pool) *App {
	// User
	userStore := user.NewStore(db)
	userService := user.NewService(userStore)
	userHandler := user.NewHandler(userService)

	// Auth
	jwt := auth.NewJwt(os.Getenv("JWT_SECRET"))
	authStore := auth.NewStore(db)
	authService := auth.NewService(authStore, jwt, userService)
	authHandler := auth.NewHandler(authService)
	authMiddleware := auth.NewMiddleware(jwt)

	// Category
	categoryStore := category.NewStore(db)
	categoryService := category.NewService(categoryStore)
	categoryHandler := category.NewHandler(categoryService)

	// Product
	productStore := product.NewStore(db)
	productService := product.NewService(productStore, categoryService)
	productHandler := product.NewHandler(productService)

	// Sale
	saleStore := sale.NewStore(db)
	saleService := sale.NewService(saleStore)
	saleHandler := sale.NewHandler(saleService)

	router := router.NewRouter(userHandler, authHandler, authMiddleware, categoryHandler, productHandler, saleHandler)

	return &App{
		Router: router,
	}
}
