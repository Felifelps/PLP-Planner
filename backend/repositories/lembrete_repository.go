package repositories

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"plp-planner/models"
)

var ErrLembreteNaoEncontrado = errors.New("lembrete não encontrado")

type LembreteRepository struct {
	db *pgxpool.Pool
}

func NewLembreteRepository(db *pgxpool.Pool) *LembreteRepository {
	return &LembreteRepository{
		db: db,
	}
}

func (r *LembreteRepository) Salvar(
	ctx context.Context,
	lembrete *models.Lembrete,
) error {

	query := `
		INSERT INTO lembretes (
			descricao,
			tipo,
			data,
			horario,
			recorrente
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	return r.db.QueryRow(
		ctx,
		query,
		lembrete.Descricao,
		lembrete.Tipo,
		lembrete.Data,
		lembrete.Horario,
		lembrete.Recorrente,
	).Scan(&lembrete.ID)
}

func (r *LembreteRepository) BuscarTodos(
	ctx context.Context,
	dataInicio string,
	dataFim string,
) ([]models.Lembrete, error) {

	query := `
		SELECT
			id,
			descricao,
			tipo,
			TO_CHAR(data, 'YYYY-MM-DD'),
			TO_CHAR(horario, 'HH24:MI'),
			recorrente
		FROM lembretes
	`

	args := []interface{}{}
	filtros := ""

	if dataInicio != "" {
		args = append(args, dataInicio)
		filtros += "data >= $" + strconv.Itoa(len(args))
	}

	if dataFim != "" {
		if filtros != "" {
			filtros += " AND "
		}

		args = append(args, dataFim)
		filtros += "data <= $" + strconv.Itoa(len(args))
	}

	if filtros != "" {
		query += " WHERE " + filtros
	}

	query += " ORDER BY data, horario"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lembretes := make([]models.Lembrete, 0)

	for rows.Next() {
		var lembrete models.Lembrete

		err := rows.Scan(
			&lembrete.ID,
			&lembrete.Descricao,
			&lembrete.Tipo,
			&lembrete.Data,
			&lembrete.Horario,
			&lembrete.Recorrente,
		)

		if err != nil {
			return nil, err
		}

		lembretes = append(lembretes, lembrete)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lembretes, nil
}

func (r *LembreteRepository) BuscarPorID(
	ctx context.Context,
	id int64,
) (*models.Lembrete, error) {

	query := `
		SELECT
			id,
			descricao,
			tipo,
			TO_CHAR(data, 'YYYY-MM-DD'),
			TO_CHAR(horario, 'HH24:MI'),
			recorrente
		FROM lembretes
		WHERE id = $1
	`

	var lembrete models.Lembrete

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&lembrete.ID,
		&lembrete.Descricao,
		&lembrete.Tipo,
		&lembrete.Data,
		&lembrete.Horario,
		&lembrete.Recorrente,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrLembreteNaoEncontrado
	}

	if err != nil {
		return nil, err
	}

	return &lembrete, nil
}

func (r *LembreteRepository) Atualizar(
	ctx context.Context,
	lembrete *models.Lembrete,
) error {

	query := `
		UPDATE lembretes
		SET
			descricao = $1,
			tipo = $2,
			data = $3,
			horario = $4,
			recorrente = $5
		WHERE id = $6
	`

	result, err := r.db.Exec(
		ctx,
		query,
		lembrete.Descricao,
		lembrete.Tipo,
		lembrete.Data,
		lembrete.Horario,
		lembrete.Recorrente,
		lembrete.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrLembreteNaoEncontrado
	}

	return nil
}

func (r *LembreteRepository) Excluir(
	ctx context.Context,
	id int64,
) error {

	query := `
		DELETE FROM lembretes
		WHERE id = $1
	`

	result, err := r.db.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrLembreteNaoEncontrado
	}

	return nil
}