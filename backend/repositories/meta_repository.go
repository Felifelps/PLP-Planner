package repositories

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"plp-planner/models"
)

var ErrMetaNaoEncontrada = errors.New("meta não encontrada")

type MetaRepository struct {
	db *pgxpool.Pool
}

func NewMetaRepository(db *pgxpool.Pool) *MetaRepository {
	return &MetaRepository{
		db: db,
	}
}

func (r *MetaRepository) Salvar(
	ctx context.Context,
	meta *models.Meta,
) error {

	query := `
		INSERT INTO metas (
			nome,
			descricao,
			categoria_id,
			status,
			periodo,
			data_inicio,
			data_fim
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	return r.db.QueryRow(
		ctx,
		query,
		meta.Nome,
		meta.Descricao,
		meta.CategoriaID,
		meta.Status,
		meta.Periodo,
		meta.DataInicio,
		meta.DataFim,
	).Scan(&meta.ID)
}

func (r *MetaRepository) BuscarTodos(
	ctx context.Context,
	dataInicio string,
	dataFim string,
) ([]models.Meta, error) {

	query := `
		SELECT
			id,
			nome,
			descricao,
			categoria_id,
			status,
			periodo,
			data_inicio,
			data_fim
		FROM metas
	`

	args := []interface{}{}
	filtros := ""

	if dataInicio != "" {
		args = append(args, dataInicio)
		filtros += "data_inicio >= $" + strconv.Itoa(len(args))
	}

	if dataFim != "" {
		if filtros != "" {
			filtros += " AND "
		}

		args = append(args, dataFim)
		filtros += "data_fim <= $" + strconv.Itoa(len(args))
	}

	if filtros != "" {
		query += " WHERE " + filtros
	}

	query += " ORDER BY id"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metas := make([]models.Meta, 0)

	for rows.Next() {
		var meta models.Meta

		err := rows.Scan(
			&meta.ID,
			&meta.Nome,
			&meta.Descricao,
			&meta.CategoriaID,
			&meta.Status,
			&meta.Periodo,
			&meta.DataInicio,
			&meta.DataFim,
		)

		if err != nil {
			return nil, err
		}

		metas = append(metas, meta)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return metas, nil
}

func (r *MetaRepository) BuscarPorID(
	ctx context.Context,
	id int64,
) (*models.Meta, error) {

	query := `
		SELECT
			id,
			nome,
			descricao,
			categoria_id,
			status,
			periodo,
			data_inicio,
			data_fim
		FROM metas
		WHERE id = $1
	`

	var meta models.Meta

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&meta.ID,
		&meta.Nome,
		&meta.Descricao,
		&meta.CategoriaID,
		&meta.Status,
		&meta.Periodo,
		&meta.DataInicio,
		&meta.DataFim,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrMetaNaoEncontrada
	}

	if err != nil {
		return nil, err
	}

	return &meta, nil
}

func (r *MetaRepository) Atualizar(
	ctx context.Context,
	meta *models.Meta,
) error {

	query := `
		UPDATE metas
		SET
			nome = $1,
			descricao = $2,
			categoria_id = $3,
			status = $4,
			periodo = $5,
			data_inicio = $6,
			data_fim = $7
		WHERE id = $8
	`

	result, err := r.db.Exec(
		ctx,
		query,
		meta.Nome,
		meta.Descricao,
		meta.CategoriaID,
		meta.Status,
		meta.Periodo,
		meta.DataInicio,
		meta.DataFim,
		meta.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrMetaNaoEncontrada
	}

	return nil
}

func (r *MetaRepository) AtualizarStatus(
	ctx context.Context,
	id int64,
	status models.Status,
) error {

	query := `
		UPDATE metas
		SET status = $1
		WHERE id = $2
	`

	result, err := r.db.Exec(
		ctx,
		query,
		status,
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrMetaNaoEncontrada
	}

	return nil
}

func (r *MetaRepository) Excluir(
	ctx context.Context,
	id int64,
) error {

	query := `
		DELETE FROM metas
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
		return ErrMetaNaoEncontrada
	}

	return nil
}
