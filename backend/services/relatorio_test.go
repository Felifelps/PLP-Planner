package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"plp-planner/models"
)

type fakeRelatorioMetaRepository struct {
	metasRetornadas []models.Meta
	erroBuscarTodos error
}

func (f *fakeRelatorioMetaRepository) BuscarTodos(
	ctx context.Context,
	dataInicio string,
	dataFim string,
) ([]models.Meta, error) {
	if f.erroBuscarTodos != nil {
		return nil, f.erroBuscarTodos
	}

	return f.metasRetornadas, nil
}

type fakeRelatorioTarefaRepository struct {
	tarefasRetornadas    []models.Tarefa
	erroBuscarPorPeriodo error
}

func (f *fakeRelatorioTarefaRepository) BuscarPorPeriodo(
	ctx context.Context,
	dataInicio string,
	dataFim string,
) ([]models.Tarefa, error) {
	if f.erroBuscarPorPeriodo != nil {
		return nil, f.erroBuscarPorPeriodo
	}

	return f.tarefasRetornadas, nil
}

func horario(s string) *string {
	return &s
}

func turno(s string) *string {
	return &s
}

func TestRelatorioServiceGerarPeriodoObrigatorio(t *testing.T) {
	service := NewRelatorioService(
		&fakeRelatorioMetaRepository{},
		&fakeRelatorioTarefaRepository{},
	)

	_, err := service.Gerar(context.Background(), "", "2026-09-30")

	if !errors.Is(err, ErrPeriodoRelatorioObrigatorio) {
		t.Fatalf("esperava ErrPeriodoRelatorioObrigatorio, obteve %v", err)
	}
}

func TestRelatorioServiceGerarDataInvalida(t *testing.T) {
	service := NewRelatorioService(
		&fakeRelatorioMetaRepository{},
		&fakeRelatorioTarefaRepository{},
	)

	_, err := service.Gerar(context.Background(), "01-09-2026", "2026-09-30")

	if !errors.Is(err, ErrPeriodoRelatorioInvalido) {
		t.Fatalf("esperava ErrPeriodoRelatorioInvalido, obteve %v", err)
	}
}

func TestRelatorioServiceGerarInicioAposFim(t *testing.T) {
	service := NewRelatorioService(
		&fakeRelatorioMetaRepository{},
		&fakeRelatorioTarefaRepository{},
	)

	_, err := service.Gerar(context.Background(), "2026-09-30", "2026-09-01")

	if !errors.Is(err, ErrPeriodoRelatorioInvalido) {
		t.Fatalf("esperava ErrPeriodoRelatorioInvalido, obteve %v", err)
	}
}

func TestRelatorioServiceGerarErroRepositorioMeta(t *testing.T) {
	erro := errors.New("falha no banco")

	service := NewRelatorioService(
		&fakeRelatorioMetaRepository{erroBuscarTodos: erro},
		&fakeRelatorioTarefaRepository{},
	)

	_, err := service.Gerar(context.Background(), "2026-09-01", "2026-09-30")

	if !errors.Is(err, erro) {
		t.Fatalf("esperava erro do repositório, obteve %v", err)
	}
}

func TestRelatorioServiceGerarErroRepositorioTarefa(t *testing.T) {
	erro := errors.New("falha no banco")

	service := NewRelatorioService(
		&fakeRelatorioMetaRepository{},
		&fakeRelatorioTarefaRepository{erroBuscarPorPeriodo: erro},
	)

	_, err := service.Gerar(context.Background(), "2026-09-01", "2026-09-30")

	if !errors.Is(err, erro) {
		t.Fatalf("esperava erro do repositório, obteve %v", err)
	}
}

func TestRelatorioServiceGerarPercentuaisSemDados(t *testing.T) {
	service := NewRelatorioService(
		&fakeRelatorioMetaRepository{},
		&fakeRelatorioTarefaRepository{},
	)

	relatorio, err := service.Gerar(context.Background(), "2026-09-01", "2026-09-30")

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if relatorio.PercentualMetasCumpridas != 0 {
		t.Errorf("esperava 0%%, obteve %v", relatorio.PercentualMetasCumpridas)
	}

	if relatorio.PercentualTarefasExecutadas != 0 {
		t.Errorf("esperava 0%%, obteve %v", relatorio.PercentualTarefasExecutadas)
	}

	if relatorio.SemanaMaisProdutiva != nil {
		t.Errorf("esperava nil, obteve %v", relatorio.SemanaMaisProdutiva)
	}

	if relatorio.TurnoMaisProdutivo != nil {
		t.Errorf("esperava nil, obteve %v", relatorio.TurnoMaisProdutivo)
	}

	if len(relatorio.CategoriasMaisRealizadasTarefas) != 0 {
		t.Errorf("esperava lista vazia, obteve %v", relatorio.CategoriasMaisRealizadasTarefas)
	}
}

func TestRelatorioServiceGerarPercentuaisComDados(t *testing.T) {
	metas := []models.Meta{
		{ID: 1, CategoriaID: 1, Status: models.StatusCumprida},
		{ID: 2, CategoriaID: 2, Status: models.StatusCumprida},
		{ID: 3, CategoriaID: 1, Status: models.StatusNaoCumprida},
		{ID: 4, CategoriaID: 1, Status: models.StatusParcialmenteCumprida},
	}

	tarefas := []models.Tarefa{
		{ID: 1, CategoriaID: 1, Status: models.StatusExecutada},
		{ID: 2, CategoriaID: 1, Status: models.StatusExecutada},
		{ID: 3, CategoriaID: 2, Status: models.StatusCancelada},
	}

	service := NewRelatorioService(
		&fakeRelatorioMetaRepository{metasRetornadas: metas},
		&fakeRelatorioTarefaRepository{tarefasRetornadas: tarefas},
	)

	relatorio, err := service.Gerar(context.Background(), "2026-09-01", "2026-09-30")

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if relatorio.TotalMetas != 4 {
		t.Errorf("esperava 4 metas, obteve %d", relatorio.TotalMetas)
	}

	if relatorio.PercentualMetasCumpridas != 50 {
		t.Errorf("esperava 50%%, obteve %v", relatorio.PercentualMetasCumpridas)
	}

	esperadoTarefas := float64(2) / float64(3) * 100
	if relatorio.PercentualTarefasExecutadas != esperadoTarefas {
		t.Errorf("esperava %v%%, obteve %v", esperadoTarefas, relatorio.PercentualTarefasExecutadas)
	}

	if len(relatorio.CategoriasMaisRealizadasTarefas) != 1 {
		t.Fatalf("esperava 1 categoria com tarefas executadas, obteve %d", len(relatorio.CategoriasMaisRealizadasTarefas))
	}

	if relatorio.CategoriasMaisRealizadasTarefas[0].CategoriaID != 1 ||
		relatorio.CategoriasMaisRealizadasTarefas[0].Total != 2 {
		t.Errorf("categoria mais realizada incorreta: %+v", relatorio.CategoriasMaisRealizadasTarefas[0])
	}

	if len(relatorio.CategoriasMaisRealizadasMetas) != 2 {
		t.Fatalf("esperava 2 categorias com metas cumpridas, obteve %d", len(relatorio.CategoriasMaisRealizadasMetas))
	}
}

func TestRelatorioServiceSemanaEMesMaisProdutivos(t *testing.T) {
	tarefas := []models.Tarefa{
		{ID: 1, CategoriaID: 1, Status: models.StatusExecutada, Data: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)},
		{ID: 2, CategoriaID: 1, Status: models.StatusExecutada, Data: time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)},
		{ID: 3, CategoriaID: 1, Status: models.StatusExecutada, Data: time.Date(2026, 9, 15, 0, 0, 0, 0, time.UTC)},
		{ID: 4, CategoriaID: 1, Status: models.StatusAdiada, Data: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)},
	}

	service := NewRelatorioService(
		&fakeRelatorioMetaRepository{},
		&fakeRelatorioTarefaRepository{tarefasRetornadas: tarefas},
	)

	relatorio, err := service.Gerar(context.Background(), "2026-09-01", "2026-09-30")

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if relatorio.SemanaMaisProdutiva == nil {
		t.Fatal("esperava semana mais produtiva calculada")
	}

	if relatorio.SemanaMaisProdutiva.Total != 2 {
		t.Errorf("esperava 2 tarefas na semana mais produtiva, obteve %d", relatorio.SemanaMaisProdutiva.Total)
	}

	if relatorio.MesMaisProdutivo == nil || relatorio.MesMaisProdutivo.Rotulo != "2026-09" {
		t.Errorf("esperava mês 2026-09 como mais produtivo, obteve %+v", relatorio.MesMaisProdutivo)
	}

	if relatorio.MesMaisProdutivo.Total != 3 {
		t.Errorf("esperava 3 tarefas no mês, obteve %d", relatorio.MesMaisProdutivo.Total)
	}
}

func TestRelatorioServiceTurnoMaisProdutivo(t *testing.T) {
	tarefas := []models.Tarefa{
		{ID: 1, CategoriaID: 1, Status: models.StatusExecutada, Turno: turno(models.TurnoManha)},
		{ID: 2, CategoriaID: 1, Status: models.StatusExecutada, HorarioInicio: horario("08:00")},
		{ID: 3, CategoriaID: 1, Status: models.StatusExecutada, HorarioInicio: horario("14:00")},
	}

	service := NewRelatorioService(
		&fakeRelatorioMetaRepository{},
		&fakeRelatorioTarefaRepository{tarefasRetornadas: tarefas},
	)

	relatorio, err := service.Gerar(context.Background(), "2026-09-01", "2026-09-30")

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if relatorio.TurnoMaisProdutivo == nil {
		t.Fatal("esperava turno mais produtivo calculado")
	}

	if relatorio.TurnoMaisProdutivo.Turno != models.TurnoManha {
		t.Errorf("esperava turno manhã, obteve %s", relatorio.TurnoMaisProdutivo.Turno)
	}

	if relatorio.TurnoMaisProdutivo.Total != 2 {
		t.Errorf("esperava 2 tarefas no turno manhã (uma explícita + uma por horário), obteve %d", relatorio.TurnoMaisProdutivo.Total)
	}
}
