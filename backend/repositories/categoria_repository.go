package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"plp-planner/models"
)

var (
	ErrCategoriaNaoEncontrada = errors.New("categoria não encontrada")
	ErrCategoriaEmUso         = errors.New("categoria em uso por tarefas ou metas")
	ErrCategoriaNomeDuplicado = errors.New("já existe uma categoria com esse nome")
)

func traduzErroCategoria(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case "23505":
		return ErrCategoriaNomeDuplicado
	case "23503":
		return ErrCategoriaEmUso
	default:
		return err
	}
}

type CategoriaRepository struct {
	db *pgxpool.Pool
}

func NewCategoriaRepository(db *pgxpool.Pool) *CategoriaRepository {
	return &CategoriaRepository{
		db: db,
	}
}

func (r *CategoriaRepository) Salvar(
	ctx context.Context,
	categoria *models.Categoria,
) error {

	query := `
		INSERT INTO categorias (nome, cor)
		VALUES ($1, $2)
		RETURNING id
	`

	err := r.db.QueryRow(
		ctx,
		query,
		categoria.Nome,
		categoria.Cor,
	).Scan(&categoria.ID)

	return traduzErroCategoria(err)
}

func (r *CategoriaRepository) BuscarTodos(
	ctx context.Context,
) ([]models.Categoria, error) {

	query := `
		SELECT id, nome, cor
		FROM categorias
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categorias := make([]models.Categoria, 0)

	for rows.Next() {
		var categoria models.Categoria

		err := rows.Scan(
			&categoria.ID,
			&categoria.Nome,
			&categoria.Cor,
		)

		if err != nil {
			return nil, err
		}

		categorias = append(categorias, categoria)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categorias, nil
}

func (r *CategoriaRepository) BuscarPorID(
	ctx context.Context,
	id int64,
) (*models.Categoria, error) {

	query := `
		SELECT id, nome, cor
		FROM categorias
		WHERE id = $1
	`

	var categoria models.Categoria

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&categoria.ID,
		&categoria.Nome,
		&categoria.Cor,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCategoriaNaoEncontrada
	}

	if err != nil {
		return nil, err
	}

	return &categoria, nil
}

func (r *CategoriaRepository) Atualizar(
	ctx context.Context,
	categoria *models.Categoria,
) error {

	query := `
		UPDATE categorias
		SET nome = $1, cor = $2
		WHERE id = $3
	`

	result, err := r.db.Exec(
		ctx,
		query,
		categoria.Nome,
		categoria.Cor,
		categoria.ID,
	)

	if err != nil {
		return traduzErroCategoria(err)
	}

	if result.RowsAffected() == 0 {
		return ErrCategoriaNaoEncontrada
	}

	return nil
}

func (r *CategoriaRepository) Excluir(
	ctx context.Context,
	id int64,
) error {

	query := `
		DELETE FROM categorias
		WHERE id = $1
	`

	result, err := r.db.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		return traduzErroCategoria(err)
	}

	if result.RowsAffected() == 0 {
		return ErrCategoriaNaoEncontrada
	}

	return nil
}
