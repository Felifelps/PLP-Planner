package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// RunMigrations executa todas as migrations pendentes.
func RunMigrations(
	ctx context.Context,
) error {
	migrationsFS, err := fs.Sub(
		migrationFiles,
		"migrations",
	)
	if err != nil {
		return fmt.Errorf(
			"erro ao carregar migrations: %w",
			err,
		)
	}

	db, err := sql.Open(
		"pgx",
		ConnectionURL(),
	)
	if err != nil {
		return fmt.Errorf(
			"erro ao abrir conexão das migrations: %w",
			err,
		)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf(
			"erro ao conectar para executar migrations: %w",
			err,
		)
	}

	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf(
			"erro ao configurar Goose: %w",
			err,
		)
	}

	if err := goose.UpContext(ctx, db, "."); err != nil {
		return fmt.Errorf(
			"erro ao executar migrations: %w",
			err,
		)
	}

	return nil
}