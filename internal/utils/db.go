package utils

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	DB() *pgxpool.Pool
}

func ExecuteTx(ctx context.Context, store Store, fn func(pgx.Tx) error) error {
	tx, err := store.DB().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
