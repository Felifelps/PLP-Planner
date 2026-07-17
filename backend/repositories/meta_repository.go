package repositories

import (
	"context"

	"fmt"

	"time"

	"plp-planner/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
			data_inicio,
			data_fim
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`

	return r.db.QueryRow(
		ctx,
		query,
		meta.Nome,
		meta.Descricao,
		meta.CategoriaID,
		meta.Status,
		meta.DataInicio,
		meta.DataFim,
	).Scan(&meta.ID)
}

func (r *MetaRepository) BuscarTodos(
	ctx context.Context,
) ([]models.Meta, error) {

	query := `
		SELECT
			id,
			nome,
			descricao,
			categoria_id,
			status,
			data_inicio,
			data_fim
		FROM metas
		ORDER BY id;
	`

	rows, err := r.db.Query(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metas := make([]models.Meta, 0)

	for rows.Next() {

		var meta models.Meta

		if err := rows.Scan(
			&meta.ID,
			&meta.Nome,
			&meta.Descricao,
			&meta.CategoriaID,
			&meta.Status,
			&meta.DataInicio,
			&meta.DataFim,
		); err != nil {
			return nil, err
		}

		metas = append(
			metas,
			meta,
		)
	}

	return metas, rows.Err()
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
			data_inicio,
			data_fim
		FROM metas
		WHERE id = $1;
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
		&meta.DataInicio,
		&meta.DataFim,
	)

	if err != nil {
		return nil, err
	}

	return &meta, nil
}

func (r *MetaRepository) BuscarPorNome(
	ctx context.Context,
	nome string,
) ([]models.Meta, error) {

	query := `
		SELECT
			id,
			nome,
			descricao,
			categoria_id,
			status,
			data_inicio,
			data_fim
		FROM metas
		WHERE nome ILIKE '%' || $1 || '%'
		ORDER BY nome;
	`

	rows, err := r.db.Query(
		ctx,
		query,
		nome,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metas := make([]models.Meta, 0)

	for rows.Next() {

		var meta models.Meta

		if err := rows.Scan(
			&meta.ID,
			&meta.Nome,
			&meta.Descricao,
			&meta.CategoriaID,
			&meta.Status,
			&meta.DataInicio,
			&meta.DataFim,
		); err != nil {
			return nil, err
		}

		metas = append(
			metas,
			meta,
		)
	}

	return metas, rows.Err()
}

func (r *MetaRepository) BuscarPorCategoria(
	ctx context.Context,
	categoriaID int64,
) ([]models.Meta, error) {

	query := `
		SELECT
			id,
			nome,
			descricao,
			categoria_id,
			status,
			data_inicio,
			data_fim
		FROM metas
		WHERE categoria_id = $1
		ORDER BY nome;
	`

	rows, err := r.db.Query(
		ctx,
		query,
		categoriaID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metas := make([]models.Meta, 0)

	for rows.Next() {

		var meta models.Meta

		if err := rows.Scan(
			&meta.ID,
			&meta.Nome,
			&meta.Descricao,
			&meta.CategoriaID,
			&meta.Status,
			&meta.DataInicio,
			&meta.DataFim,
		); err != nil {
			return nil, err
		}

		metas = append(
			metas,
			meta,
		)
	}

	return metas, rows.Err()
}

func (r *MetaRepository) BuscarPorStatus(
	ctx context.Context,
	status string,
) ([]models.Meta, error) {

	query := `
		SELECT
			id,
			nome,
			descricao,
			categoria_id,
			status,
			data_inicio,
			data_fim
		FROM metas
		WHERE status = $1
		ORDER BY data_fim;
	`

	rows, err := r.db.Query(
		ctx,
		query,
		status,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metas := make([]models.Meta, 0)

	for rows.Next() {

		var meta models.Meta

		if err := rows.Scan(
			&meta.ID,
			&meta.Nome,
			&meta.Descricao,
			&meta.CategoriaID,
			&meta.Status,
			&meta.DataInicio,
			&meta.DataFim,
		); err != nil {
			return nil, err
		}

		metas = append(
			metas,
			meta,
		)
	}

	return metas, rows.Err()
}

func (r *MetaRepository) BuscarPorPeriodo(
	ctx context.Context,
	dataInicio time.Time,
	dataFim time.Time,
) ([]models.Meta, error) {

	query := `
		SELECT
			id,
			nome,
			descricao,
			categoria_id,
			status,
			data_inicio,
			data_fim
		FROM metas
		WHERE data_inicio >= $1
		  AND data_fim <= $2
		ORDER BY data_inicio;
	`

	rows, err := r.db.Query(
		ctx,
		query,
		dataInicio,
		dataFim,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metas := make([]models.Meta, 0)

	for rows.Next() {

		var meta models.Meta

		if err := rows.Scan(
			&meta.ID,
			&meta.Nome,
			&meta.Descricao,
			&meta.CategoriaID,
			&meta.Status,
			&meta.DataInicio,
			&meta.DataFim,
		); err != nil {
			return nil, err
		}

		metas = append(
			metas,
			meta,
		)
	}

	return metas, rows.Err()
}

func (r *MetaRepository) BuscarRecentes(
	ctx context.Context,
	limite int,
) ([]models.Meta, error) {

	query := `
		SELECT
			id,
			nome,
			descricao,
			categoria_id,
			status,
			data_inicio,
			data_fim
		FROM metas
		ORDER BY id DESC
		LIMIT $1;
	`

	rows, err := r.db.Query(
		ctx,
		query,
		limite,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metas := make([]models.Meta, 0)

	for rows.Next() {

		var meta models.Meta

		if err := rows.Scan(
			&meta.ID,
			&meta.Nome,
			&meta.Descricao,
			&meta.CategoriaID,
			&meta.Status,
			&meta.DataInicio,
			&meta.DataFim,
		); err != nil {
			return nil, err
		}

		metas = append(
			metas,
			meta,
		)
	}

	return metas, rows.Err()
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
			data_inicio = $5,
			data_fim = $6
		WHERE id = $7;
	`

	result, err := r.db.Exec(
		ctx,
		query,
		meta.Nome,
		meta.Descricao,
		meta.CategoriaID,
		meta.Status,
		meta.DataInicio,
		meta.DataFim,
		meta.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("meta com id %d não encontrada", meta.ID)
	}

	return nil
}

func (r *MetaRepository) AtualizarStatus(
	ctx context.Context,
	id int64,
	status string,
) error {

	query := `
		UPDATE metas
		SET
			status = $1
		WHERE id = $2;
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

	// Verifica se a linha realmente existia
	if result.RowsAffected() == 0 {
		return fmt.Errorf("meta com id %d não encontrada", id)
	}

	return nil
}

func (r *MetaRepository) Excluir(
	ctx context.Context,
	id int64,
) error {

	query := `
		DELETE FROM metas
		WHERE id = $1;
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
		return fmt.Errorf("meta com id %d não encontrada", id)
	}

	return nil
}
