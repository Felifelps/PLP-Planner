package repositories

import (
	"context"

	"plp-planner/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ExemploRepository struct {
	db *pgxpool.Pool
}

func NewExemploRepository(db *pgxpool.Pool) *ExemploRepository {
	return &ExemploRepository{db: db}
}

func (r *ExemploRepository) BuscarTodos(
	ctx context.Context,
) ([]models.Exemplo, error) {
	rows, err := r.db.Query(
		ctx,
		"SELECT id, nome FROM exemplos",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exemplos := make([]models.Exemplo, 0)

	for rows.Next() {
		var exemplo models.Exemplo

		if err := rows.Scan(
			&exemplo.ID,
			&exemplo.Nome,
		); err != nil {
			return nil, err
		}

		exemplos = append(exemplos, exemplo)
	}

	return exemplos, rows.Err()
}