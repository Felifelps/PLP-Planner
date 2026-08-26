package services

import (
	"context"
	"errors"
	"strings"

	"plp-planner/models"
)

var (
	ErrTarefaInvalida          = errors.New("tarefa inválida")
	ErrTarefaIDInvalido        = errors.New("id inválido")
	ErrDescricaoObrigatoria    = errors.New("descrição obrigatória")
	ErrTarefaCategoriaInvalida = errors.New("categoria inválida")
	ErrTarefaStatusInvalido    = errors.New("status inválido")
)

type TarefaRepository interface {
	Salvar(
		ctx context.Context,
		tarefa *models.Tarefa,
	) error

	BuscarTodos(
		ctx context.Context,
		data string,
		categoriaID string,
	) ([]models.Tarefa, error)

	BuscarPorID(
		ctx context.Context,
		id int64,
	) (*models.Tarefa, error)

	Atualizar(
		ctx context.Context,
		tarefa *models.Tarefa,
	) error

	AtualizarStatus(
		ctx context.Context,
		id int64,
		status models.StatusTarefa,
	) error

	Excluir(
		ctx context.Context,
		id int64,
	) error
}

type TarefaService struct {
	repository TarefaRepository
}

func NewTarefaService(
	repository TarefaRepository,
) *TarefaService {
	return &TarefaService{
		repository: repository,
	}
}

func (s *TarefaService) validarTarefa(
	tarefa *models.Tarefa,
) error {

	if tarefa == nil {
		return ErrTarefaInvalida
	}

	if strings.TrimSpace(tarefa.Descricao) == "" {
		return ErrDescricaoObrigatoria
	}

	if tarefa.CategoriaID <= 0 {
		return ErrTarefaCategoriaInvalida
	}

	if err := tarefa.Validate(); err != nil {
		return err
	}

	return nil
}

func (s *TarefaService) Salvar(
	ctx context.Context,
	tarefa *models.Tarefa,
) error {

	if err := s.validarTarefa(tarefa); err != nil {
		return err
	}

	return s.repository.Salvar(
		ctx,
		tarefa,
	)
}

func (s *TarefaService) BuscarTodos(
	ctx context.Context,
	data string,
	categoriaID string,
) ([]models.Tarefa, error) {

	return s.repository.BuscarTodos(
		ctx,
		data,
		categoriaID,
	)
}

func (s *TarefaService) BuscarPorID(
	ctx context.Context,
	id int64,
) (*models.Tarefa, error) {

	if id <= 0 {
		return nil, ErrTarefaIDInvalido
	}

	return s.repository.BuscarPorID(
		ctx,
		id,
	)
}

func (s *TarefaService) Atualizar(
	ctx context.Context,
	tarefa *models.Tarefa,
) error {

	if tarefa == nil {
		return ErrTarefaInvalida
	}

	if tarefa.ID <= 0 {
		return ErrTarefaIDInvalido
	}

	if err := s.validarTarefa(tarefa); err != nil {
		return err
	}

	return s.repository.Atualizar(
		ctx,
		tarefa,
	)
}

func (s *TarefaService) AtualizarStatus(
	ctx context.Context,
	id int64,
	status models.StatusTarefa,
) error {

	if id <= 0 {
		return ErrTarefaIDInvalido
	}

	if !models.StatusTarefaValido(status) {
		return ErrTarefaStatusInvalido
	}

	return s.repository.AtualizarStatus(
		ctx,
		id,
		status,
	)
}

func (s *TarefaService) Excluir(
	ctx context.Context,
	id int64,
) error {

	if id <= 0 {
		return ErrTarefaIDInvalido
	}

	return s.repository.Excluir(
		ctx,
		id,
	)
}
