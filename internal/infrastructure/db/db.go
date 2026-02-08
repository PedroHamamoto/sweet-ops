package db

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := "postgres://" + os.Getenv("DATABASE_USER") + ":" + os.Getenv("DATABASE_PASSWORD") + "@" + os.Getenv("DATABASE_HOST") + ":" + os.Getenv("DATABASE_PORT") + "/" + os.Getenv("DATABASE_NAME") + "?sslmode=" + os.Getenv("DATABASE_SSLMODE")
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour

	return pgxpool.NewWithConfig(ctx, cfg)
}
