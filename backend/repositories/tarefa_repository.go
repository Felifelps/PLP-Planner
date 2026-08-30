package repositories

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"plp-planner/models"
)

var ErrTarefaNaoEncontrada = errors.New("tarefa não encontrada")

type TarefaRepository struct {
	db *pgxpool.Pool
}

func NewTarefaRepository(db *pgxpool.Pool) *TarefaRepository {
	return &TarefaRepository{
		db: db,
	}
}

func (r *TarefaRepository) Salvar(
	ctx context.Context,
	tarefa *models.Tarefa,
) error {

	query := `
		INSERT INTO tarefas (
			descricao,
			categoria_id,
			data,
			horario_inicio,
			duracao,
			turno,
			status,
			prioridade
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	return r.db.QueryRow(
		ctx,
		query,
		tarefa.Descricao,
		tarefa.CategoriaID,
		tarefa.Data,
		tarefa.HorarioInicio,
		tarefa.Duracao,
		tarefa.Turno,
		tarefa.Status,
		tarefa.Prioridade,
	).Scan(&tarefa.ID)
}

func (r *TarefaRepository) BuscarTodos(
	ctx context.Context,
	data string,
	categoriaID string,
) ([]models.Tarefa, error) {

	query := `
		SELECT
			id,
			descricao,
			categoria_id,
			data,
			horario_inicio,
			duracao,
			turno,
			status,
			prioridade
		FROM tarefas
	`

	args := []interface{}{}
	filtros := ""

	if data != "" {
		args = append(args, data)
		filtros += "data = $" + strconv.Itoa(len(args))
	}

	if categoriaID != "" {
		if filtros != "" {
			filtros += " AND "
		}

		args = append(args, categoriaID)
		filtros += "categoria_id = $" + strconv.Itoa(len(args))
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

	tarefas := make([]models.Tarefa, 0)

	for rows.Next() {
		var tarefa models.Tarefa

		err := rows.Scan(
			&tarefa.ID,
			&tarefa.Descricao,
			&tarefa.CategoriaID,
			&tarefa.Data,
			&tarefa.HorarioInicio,
			&tarefa.Duracao,
			&tarefa.Turno,
			&tarefa.Status,
			&tarefa.Prioridade,
		)

		if err != nil {
			return nil, err
		}

		tarefas = append(tarefas, tarefa)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tarefas, nil
}

func (r *TarefaRepository) BuscarPorPeriodo(
	ctx context.Context,
	dataInicio string,
	dataFim string,
) ([]models.Tarefa, error) {

	query := `
		SELECT
			id,
			descricao,
			categoria_id,
			data,
			horario_inicio,
			duracao,
			turno,
			status,
			prioridade
		FROM tarefas
		WHERE data >= $1 AND data <= $2
		ORDER BY data, id
	`

	rows, err := r.db.Query(ctx, query, dataInicio, dataFim)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tarefas := make([]models.Tarefa, 0)

	for rows.Next() {
		var tarefa models.Tarefa

		err := rows.Scan(
			&tarefa.ID,
			&tarefa.Descricao,
			&tarefa.CategoriaID,
			&tarefa.Data,
			&tarefa.HorarioInicio,
			&tarefa.Duracao,
			&tarefa.Turno,
			&tarefa.Status,
			&tarefa.Prioridade,
		)

		if err != nil {
			return nil, err
		}

		tarefas = append(tarefas, tarefa)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tarefas, nil
}

func (r *TarefaRepository) BuscarPorID(
	ctx context.Context,
	id int64,
) (*models.Tarefa, error) {

	query := `
		SELECT
			id,
			descricao,
			categoria_id,
			data,
			horario_inicio,
			duracao,
			turno,
			status,
			prioridade
		FROM tarefas
		WHERE id = $1
	`

	var tarefa models.Tarefa

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&tarefa.ID,
		&tarefa.Descricao,
		&tarefa.CategoriaID,
		&tarefa.Data,
		&tarefa.HorarioInicio,
		&tarefa.Duracao,
		&tarefa.Turno,
		&tarefa.Status,
		&tarefa.Prioridade,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTarefaNaoEncontrada
	}

	if err != nil {
		return nil, err
	}

	return &tarefa, nil
}

func (r *TarefaRepository) Atualizar(
	ctx context.Context,
	tarefa *models.Tarefa,
) error {

	query := `
		UPDATE tarefas
		SET
			descricao = $1,
			categoria_id = $2,
			data = $3,
			horario_inicio = $4,
			duracao = $5,
			turno = $6,
			status = $7,
			prioridade = $8
		WHERE id = $9
	`

	result, err := r.db.Exec(
		ctx,
		query,
		tarefa.Descricao,
		tarefa.CategoriaID,
		tarefa.Data,
		tarefa.HorarioInicio,
		tarefa.Duracao,
		tarefa.Turno,
		tarefa.Status,
		tarefa.Prioridade,
		tarefa.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrTarefaNaoEncontrada
	}

	return nil
}

func (r *TarefaRepository) AtualizarStatus(
	ctx context.Context,
	id int64,
	status models.StatusTarefa,
) error {

	query := `
		UPDATE tarefas
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
		return ErrTarefaNaoEncontrada
	}

	return nil
}

func (r *TarefaRepository) Excluir(
	ctx context.Context,
	id int64,
) error {

	query := `
		DELETE FROM tarefas
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
		return ErrTarefaNaoEncontrada
	}

	return nil
}
