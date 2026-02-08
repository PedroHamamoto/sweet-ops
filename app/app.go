package app

import (
	"net/http"
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

	router := router.NewRouter(userHandler)

	return &App{
		Router: router,
	}
}
