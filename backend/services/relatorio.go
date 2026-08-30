package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"plp-planner/models"
)

var (
	ErrPeriodoRelatorioObrigatorio = errors.New("data_inicio e data_fim são obrigatórias")
	ErrPeriodoRelatorioInvalido    = errors.New("período inválido")
)

type RelatorioMetaRepository interface {
	BuscarTodos(
		ctx context.Context,
		dataInicio string,
		dataFim string,
	) ([]models.Meta, error)
}

type RelatorioTarefaRepository interface {
	BuscarPorPeriodo(
		ctx context.Context,
		dataInicio string,
		dataFim string,
	) ([]models.Tarefa, error)
}

type RelatorioService struct {
	metaRepository   RelatorioMetaRepository
	tarefaRepository RelatorioTarefaRepository
}

func NewRelatorioService(
	metaRepository RelatorioMetaRepository,
	tarefaRepository RelatorioTarefaRepository,
) *RelatorioService {
	return &RelatorioService{
		metaRepository:   metaRepository,
		tarefaRepository: tarefaRepository,
	}
}

func (s *RelatorioService) Gerar(
	ctx context.Context,
	dataInicio string,
	dataFim string,
) (*models.Relatorio, error) {

	if dataInicio == "" || dataFim == "" {
		return nil, ErrPeriodoRelatorioObrigatorio
	}

	inicio, err := time.Parse("2006-01-02", dataInicio)
	if err != nil {
		return nil, ErrPeriodoRelatorioInvalido
	}

	fim, err := time.Parse("2006-01-02", dataFim)
	if err != nil {
		return nil, ErrPeriodoRelatorioInvalido
	}

	if inicio.After(fim) {
		return nil, ErrPeriodoRelatorioInvalido
	}

	metas, err := s.metaRepository.BuscarTodos(ctx, dataInicio, dataFim)
	if err != nil {
		return nil, err
	}

	tarefas, err := s.tarefaRepository.BuscarPorPeriodo(ctx, dataInicio, dataFim)
	if err != nil {
		return nil, err
	}

	return &models.Relatorio{
		DataInicio:   dataInicio,
		DataFim:      dataFim,
		TotalMetas:   len(metas),
		TotalTarefas: len(tarefas),

		PercentualMetasCumpridas:    percentual(contarMetasCumpridas(metas), len(metas)),
		PercentualTarefasExecutadas: percentual(contarTarefasExecutadas(tarefas), len(tarefas)),

		SemanaMaisProdutiva: periodoMaisProdutivo(tarefas, rotuloSemana),
		MesMaisProdutivo:    periodoMaisProdutivo(tarefas, rotuloMes),
		TurnoMaisProdutivo:  turnoMaisProdutivo(tarefas),

		CategoriasMaisRealizadasTarefas: categoriasMaisRealizadasTarefas(tarefas),
		CategoriasMaisRealizadasMetas:   categoriasMaisRealizadasMetas(metas),
	}, nil
}

func percentual(realizadas int, total int) float64 {
	if total == 0 {
		return 0
	}

	return float64(realizadas) / float64(total) * 100
}

func contarMetasCumpridas(metas []models.Meta) int {
	total := 0

	for _, meta := range metas {
		if meta.Status == models.StatusCumprida {
			total++
		}
	}

	return total
}

func contarTarefasExecutadas(tarefas []models.Tarefa) int {
	total := 0

	for _, tarefa := range tarefas {
		if tarefa.Status == models.StatusExecutada {
			total++
		}
	}

	return total
}

func rotuloSemana(data time.Time) string {
	ano, semana := data.ISOWeek()
	return fmt.Sprintf("%d-W%02d", ano, semana)
}

func rotuloMes(data time.Time) string {
	return data.Format("2006-01")
}

func periodoMaisProdutivo(
	tarefas []models.Tarefa,
	rotulo func(time.Time) string,
) *models.PeriodoContagem {

	contagem := map[string]int{}

	for _, tarefa := range tarefas {
		if tarefa.Status != models.StatusExecutada {
			continue
		}

		contagem[rotulo(tarefa.Data)]++
	}

	if len(contagem) == 0 {
		return nil
	}

	melhorRotulo := ""
	melhorTotal := -1

	for r, total := range contagem {
		if total > melhorTotal ||
			(total == melhorTotal && r < melhorRotulo) {

			melhorRotulo = r
			melhorTotal = total
		}
	}

	return &models.PeriodoContagem{
		Rotulo: melhorRotulo,
		Total:  melhorTotal,
	}
}

func turnoDaTarefa(tarefa models.Tarefa) string {
	if tarefa.Turno != nil {
		return *tarefa.Turno
	}

	if tarefa.HorarioInicio == nil {
		return ""
	}

	horario := *tarefa.HorarioInicio

	switch {
	case horario >= "06:00" && horario < "12:00":
		return models.TurnoManha
	case horario >= "12:00" && horario < "18:00":
		return models.TurnoTarde
	default:
		return models.TurnoNoite
	}
}

func turnoMaisProdutivo(tarefas []models.Tarefa) *models.TurnoContagem {
	contagem := map[string]int{}

	for _, tarefa := range tarefas {
		if tarefa.Status != models.StatusExecutada {
			continue
		}

		turno := turnoDaTarefa(tarefa)
		if turno == "" {
			continue
		}

		contagem[turno]++
	}

	if len(contagem) == 0 {
		return nil
	}

	melhorTurno := ""
	melhorTotal := -1

	for turno, total := range contagem {
		if total > melhorTotal ||
			(total == melhorTotal && turno < melhorTurno) {

			melhorTurno = turno
			melhorTotal = total
		}
	}

	return &models.TurnoContagem{
		Turno: melhorTurno,
		Total: melhorTotal,
	}
}

func categoriasMaisRealizadasTarefas(tarefas []models.Tarefa) []models.CategoriaContagem {
	contagem := map[int64]int{}

	for _, tarefa := range tarefas {
		if tarefa.Status != models.StatusExecutada {
			continue
		}

		contagem[tarefa.CategoriaID]++
	}

	return ordenarCategorias(contagem)
}

func categoriasMaisRealizadasMetas(metas []models.Meta) []models.CategoriaContagem {
	contagem := map[int64]int{}

	for _, meta := range metas {
		if meta.Status != models.StatusCumprida {
			continue
		}

		contagem[meta.CategoriaID]++
	}

	return ordenarCategorias(contagem)
}

func ordenarCategorias(contagem map[int64]int) []models.CategoriaContagem {
	resultado := make([]models.CategoriaContagem, 0, len(contagem))

	for categoriaID, total := range contagem {
		resultado = append(resultado, models.CategoriaContagem{
			CategoriaID: categoriaID,
			Total:       total,
		})
	}

	sort.Slice(resultado, func(i, j int) bool {
		if resultado[i].Total == resultado[j].Total {
			return resultado[i].CategoriaID < resultado[j].CategoriaID
		}

		return resultado[i].Total > resultado[j].Total
	})

	return resultado
}
