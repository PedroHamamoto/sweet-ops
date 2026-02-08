package app

import (
	"net/http"
	"os"
	"sweet-ops/internal/auth"
	"sweet-ops/internal/http/router"
	"sweet-ops/internal/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Router http.Handler
}

func New(db *pgxpool.Pool) *App {
	userStore := user.NewStore(db)
	userService := user.NewService(userStore)
	userHandler := user.NewHandler(userService)

	jwt := auth.NewJwt(os.Getenv("JWT_SECRET"))
	authStore := auth.NewStore(db)
	authService := auth.NewService(authStore, jwt, userService)
	authHandler := auth.NewHandler(authService)

	router := router.NewRouter(userHandler, authHandler)

	return &App{
		Router: router,
	}
}
