package services

import (
	"context"

	"plp-planner/models"
)

type ExemploRepository interface {
	BuscarTodos(ctx context.Context) ([]models.Exemplo, error)
}

type ExemploService struct {
	repository ExemploRepository
}

func NewExemploService(
	repository ExemploRepository,
) *ExemploService {
	return &ExemploService{
		repository: repository,
	}
}

func (s *ExemploService) BuscarTodos(
	ctx context.Context,
) ([]models.Exemplo, error) {
	return s.repository.BuscarTodos(ctx)
}