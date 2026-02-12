package app

import (
	"net/http"
	"os"
	"sweet-ops/internal/auth"
	"sweet-ops/internal/category"
	"sweet-ops/internal/http/router"
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

	router := router.NewRouter(userHandler, authHandler, authMiddleware, categoryHandler)

	return &App{
		Router: router,
	}
}
