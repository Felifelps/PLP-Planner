package services

import (
	"context"
	"errors"
	"strings"

	"plp-planner/models"
)

var (
	ErrMetaInvalida      = errors.New("meta inválida")
	ErrIDInvalido        = errors.New("id inválido")
	ErrNomeObrigatorio   = errors.New("nome obrigatório")
	ErrCategoriaInvalida = errors.New("categoria inválida")
	ErrStatusInvalido    = errors.New("status inválido")
	ErrPeriodoInvalido   = errors.New("período inválido")
)

type MetaRepository interface {
	Salvar(
		ctx context.Context,
		meta *models.Meta,
	) error

	BuscarTodos(
		ctx context.Context,
	) ([]models.Meta, error)

	BuscarPorID(
		ctx context.Context,
		id int64,
	) (*models.Meta, error)

	Atualizar(
		ctx context.Context,
		meta *models.Meta,
	) error

	AtualizarStatus(
		ctx context.Context,
		id int64,
		status string,
	) error

	Excluir(
		ctx context.Context,
		id int64,
	) error
}

type MetaService struct {
	repository MetaRepository
}

func NewMetaService(
	repository MetaRepository,
) *MetaService {

	return &MetaService{
		repository: repository,
	}
}

func (s *MetaService) validarMeta(
	meta *models.Meta,
) error {

	if meta == nil {
		return ErrMetaInvalida
	}

	if strings.TrimSpace(meta.Nome) == "" {
		return ErrNomeObrigatorio
	}

	if meta.CategoriaID <= 0 {
		return ErrCategoriaInvalida
	}

	if !models.StatusValido(
		meta.Status,
	) {
		return ErrStatusInvalido
	}

	if !meta.PeriodoValido() {
		return ErrPeriodoInvalido
	}

	return nil
}

func (s *MetaService) Salvar(
	ctx context.Context,
	meta *models.Meta,
) error {

	if err := s.validarMeta(meta); err != nil {
		return err
	}

	return s.repository.Salvar(
		ctx,
		meta,
	)
}

func (s *MetaService) BuscarTodos(
	ctx context.Context,
) ([]models.Meta, error) {

	return s.repository.BuscarTodos(
		ctx,
	)
}

func (s *MetaService) BuscarPorID(
	ctx context.Context,
	id int64,
) (*models.Meta, error) {

	if id <= 0 {
		return nil, ErrIDInvalido
	}

	return s.repository.BuscarPorID(
		ctx,
		id,
	)
}

func (s *MetaService) Atualizar(
	ctx context.Context,
	meta *models.Meta,
) error {

	if meta == nil {
		return ErrMetaInvalida
	}

	if meta.ID <= 0 {
		return ErrIDInvalido
	}

	if err := s.validarMeta(meta); err != nil {
		return err
	}

	return s.repository.Atualizar(
		ctx,
		meta,
	)
}

func (s *MetaService) AtualizarStatus(
	ctx context.Context,
	id int64,
	status string,
) error {

	if id <= 0 {
		return ErrIDInvalido
	}

	if !models.StatusValido(
		status,
	) {
		return ErrStatusInvalido
	}

	return s.repository.AtualizarStatus(
		ctx,
		id,
		status,
	)
}

func (s *MetaService) Excluir(
	ctx context.Context,
	id int64,
) error {

	if id <= 0 {
		return ErrIDInvalido
	}

	return s.repository.Excluir(
		ctx,
		id,
	)
}
