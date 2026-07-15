package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectionURL cria a URL de conexão com o PostgreSQL.
func ConnectionURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

// Connect cria e valida o pool de conexões.
func Connect(
	ctx context.Context,
) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(
		ctx,
		ConnectionURL(),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"erro ao criar pool de conexões: %w",
			err,
		)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, fmt.Errorf(
			"erro ao conectar ao PostgreSQL: %w",
			err,
		)
	}

	return pool, nil
}