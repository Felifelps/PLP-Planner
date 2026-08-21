package services

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"plp-planner/models"
)

var (
	ErrLembreteInvalido             = errors.New("lembrete inválido")
	ErrDescricaoLembreteObrigatoria = errors.New("descrição obrigatória")
	ErrDataLembreteInvalida         = errors.New("data inválida")
	ErrHorarioLembreteInvalido      = errors.New("horário inválido")
	ErrPeriodoLembreteInvalido      = errors.New("período inválido")
)

type LembreteRepository interface {
	Salvar(
		ctx context.Context,
		lembrete *models.Lembrete,
	) error

	BuscarTodos(
		ctx context.Context,
		dataInicio string,
		dataFim string,
	) ([]models.Lembrete, error)

	BuscarPorID(
		ctx context.Context,
		id int64,
	) (*models.Lembrete, error)

	Atualizar(
		ctx context.Context,
		lembrete *models.Lembrete,
	) error

	Excluir(
		ctx context.Context,
		id int64,
	) error
}

type LembreteService struct {
	repository LembreteRepository
}

func NewLembreteService(
	repository LembreteRepository,
) *LembreteService {
	return &LembreteService{
		repository: repository,
	}
}

func (s *LembreteService) validarLembrete(
	lembrete *models.Lembrete,
) error {

	if lembrete == nil {
		return ErrLembreteInvalido
	}

	if strings.TrimSpace(lembrete.Descricao) == "" {
		return ErrDescricaoLembreteObrigatoria
	}

	if _, err := time.Parse(
		"2006-01-02",
		lembrete.Data,
	); err != nil {
		return ErrDataLembreteInvalida
	}

	if !horarioValido(lembrete.Horario) {
		return ErrHorarioLembreteInvalido
	}

	if err := lembrete.Validate(); err != nil {
		return err
	}

	return nil
}

func horarioValido(
	horario string,
) bool {

	if _, err := time.Parse(
		"15:04",
		horario,
	); err == nil {
		return true
	}

	if _, err := time.Parse(
		"15:04:05",
		horario,
	); err == nil {
		return true
	}

	return false
}

func (s *LembreteService) Salvar(
	ctx context.Context,
	lembrete *models.Lembrete,
) error {

	if err := s.validarLembrete(lembrete); err != nil {
		return err
	}

	return s.repository.Salvar(
		ctx,
		lembrete,
	)
}

func (s *LembreteService) BuscarTodos(
	ctx context.Context,
	dataInicio string,
	dataFim string,
) ([]models.Lembrete, error) {

	var inicio time.Time
	var fim time.Time
	var err error

	if dataInicio != "" {
		inicio, err = time.Parse(
			"2006-01-02",
			dataInicio,
		)

		if err != nil {
			return nil, ErrDataLembreteInvalida
		}
	}

	if dataFim != "" {
		fim, err = time.Parse(
			"2006-01-02",
			dataFim,
		)

		if err != nil {
			return nil, ErrDataLembreteInvalida
		}
	}

	if dataInicio != "" &&
		dataFim != "" &&
		inicio.After(fim) {

		return nil, ErrPeriodoLembreteInvalido
	}

	/*
		Quando existe um intervalo completo, buscamos
		todos os lembretes até a data final.

		Isso é necessário porque um lembrete recorrente
		pode ter sido criado antes da data inicial e ainda
		precisa aparecer nas semanas seguintes.
	*/
	if dataInicio != "" && dataFim != "" {

		lembretes, err := s.repository.BuscarTodos(
			ctx,
			"",
			dataFim,
		)

		if err != nil {
			return nil, err
		}

		return gerarOcorrenciasSemanais(
			lembretes,
			inicio,
			fim,
		)
	}

	return s.repository.BuscarTodos(
		ctx,
		dataInicio,
		dataFim,
	)
}

func gerarOcorrenciasSemanais(
	lembretes []models.Lembrete,
	inicio time.Time,
	fim time.Time,
) ([]models.Lembrete, error) {

	resultado := make([]models.Lembrete, 0)

	for _, lembrete := range lembretes {

		dataLembrete, err := time.Parse(
			"2006-01-02",
			lembrete.Data,
		)

		if err != nil {
			return nil, ErrDataLembreteInvalida
		}

		if !lembrete.Recorrente {

			if !dataLembrete.Before(inicio) &&
				!dataLembrete.After(fim) {

				resultado = append(
					resultado,
					lembrete,
				)
			}

			continue
		}

		/*
			Descobre a primeira ocorrência do lembrete
			dentro do intervalo solicitado.
		*/
		ocorrencia := dataLembrete

		for ocorrencia.Before(inicio) {
			ocorrencia = ocorrencia.AddDate(
				0,
				0,
				7,
			)
		}

		/*
			Adiciona uma ocorrência a cada 7 dias
			até ultrapassar a data final.
		*/
		for !ocorrencia.After(fim) {

			copia := lembrete

			copia.Data = ocorrencia.Format(
				"2006-01-02",
			)

			resultado = append(
				resultado,
				copia,
			)

			ocorrencia = ocorrencia.AddDate(
				0,
				0,
				7,
			)
		}
	}

	sort.Slice(
		resultado,
		func(i, j int) bool {

			if resultado[i].Data == resultado[j].Data {
				return resultado[i].Horario <
					resultado[j].Horario
			}

			return resultado[i].Data <
				resultado[j].Data
		},
	)

	return resultado, nil
}

func (s *LembreteService) BuscarPorID(
	ctx context.Context,
	id int64,
) (*models.Lembrete, error) {

	if id <= 0 {
		return nil, ErrIDInvalido
	}

	return s.repository.BuscarPorID(
		ctx,
		id,
	)
}

func (s *LembreteService) Atualizar(
	ctx context.Context,
	lembrete *models.Lembrete,
) error {

	if lembrete == nil {
		return ErrLembreteInvalido
	}

	if lembrete.ID <= 0 {
		return ErrIDInvalido
	}

	if err := s.validarLembrete(lembrete); err != nil {
		return err
	}

	return s.repository.Atualizar(
		ctx,
		lembrete,
	)
}

func (s *LembreteService) Excluir(
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