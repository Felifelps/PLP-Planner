package services

import (
	"context"
	"errors"
	"strings"

	"plp-planner/models"
)

var (
	ErrCategoriaInvalidaEntrada = errors.New("categoria inválida")
	ErrCategoriaIDInvalido      = errors.New("id inválido")
	ErrCategoriaNomeObrigatorio = errors.New("nome obrigatório")
)

type CategoriaRepository interface {
	Salvar(
		ctx context.Context,
		categoria *models.Categoria,
	) error

	BuscarTodos(
		ctx context.Context,
	) ([]models.Categoria, error)

	BuscarPorID(
		ctx context.Context,
		id int64,
	) (*models.Categoria, error)

	Atualizar(
		ctx context.Context,
		categoria *models.Categoria,
	) error

	Excluir(
		ctx context.Context,
		id int64,
	) error
}

type CategoriaService struct {
	repository CategoriaRepository
}

func NewCategoriaService(
	repository CategoriaRepository,
) *CategoriaService {
	return &CategoriaService{
		repository: repository,
	}
}

func (s *CategoriaService) validarCategoria(
	categoria *models.Categoria,
) error {

	if categoria == nil {
		return ErrCategoriaInvalidaEntrada
	}

	if strings.TrimSpace(categoria.Nome) == "" {
		return ErrCategoriaNomeObrigatorio
	}

	if err := categoria.Validate(); err != nil {
		return err
	}

	return nil
}

func (s *CategoriaService) Salvar(
	ctx context.Context,
	categoria *models.Categoria,
) error {

	if err := s.validarCategoria(categoria); err != nil {
		return err
	}

	return s.repository.Salvar(
		ctx,
		categoria,
	)
}

func (s *CategoriaService) BuscarTodos(
	ctx context.Context,
) ([]models.Categoria, error) {

	return s.repository.BuscarTodos(ctx)
}

func (s *CategoriaService) BuscarPorID(
	ctx context.Context,
	id int64,
) (*models.Categoria, error) {

	if id <= 0 {
		return nil, ErrCategoriaIDInvalido
	}

	return s.repository.BuscarPorID(
		ctx,
		id,
	)
}

func (s *CategoriaService) Atualizar(
	ctx context.Context,
	categoria *models.Categoria,
) error {

	if categoria == nil {
		return ErrCategoriaInvalidaEntrada
	}

	if categoria.ID <= 0 {
		return ErrCategoriaIDInvalido
	}

	if err := s.validarCategoria(categoria); err != nil {
		return err
	}

	return s.repository.Atualizar(
		ctx,
		categoria,
	)
}

func (s *CategoriaService) Excluir(
	ctx context.Context,
	id int64,
) error {

	if id <= 0 {
		return ErrCategoriaIDInvalido
	}

	return s.repository.Excluir(
		ctx,
		id,
	)
}
