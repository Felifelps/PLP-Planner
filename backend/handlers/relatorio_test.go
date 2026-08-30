package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"plp-planner/models"
	"plp-planner/services"
)

type fakeRelatorioMetaRepository struct {
	metas           []models.Meta
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

	return f.metas, nil
}

type fakeRelatorioTarefaRepository struct {
	tarefas              []models.Tarefa
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

	return f.tarefas, nil
}

func criarRelatorioHandlerTeste(
	metaRepository *fakeRelatorioMetaRepository,
	tarefaRepository *fakeRelatorioTarefaRepository,
) *RelatorioHandler {
	service := services.NewRelatorioService(metaRepository, tarefaRepository)

	return NewRelatorioHandler(service)
}

func TestRelatorioHandlerGerar(t *testing.T) {
	metas := []models.Meta{
		{ID: 1, CategoriaID: 1, Status: models.StatusCumprida},
	}

	tarefas := []models.Tarefa{
		{ID: 1, CategoriaID: 1, Status: models.StatusExecutada},
	}

	handler := criarRelatorioHandlerTeste(
		&fakeRelatorioMetaRepository{metas: metas},
		&fakeRelatorioTarefaRepository{tarefas: tarefas},
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/relatorios?data_inicio=2026-09-01&data_fim=2026-09-30",
		nil,
	)

	response := httptest.NewRecorder()

	handler.Gerar(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("esperava status 200, obteve %d: %s", response.Code, response.Body.String())
	}

	var relatorio models.Relatorio

	if err := json.NewDecoder(response.Body).Decode(&relatorio); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}

	if relatorio.TotalMetas != 1 || relatorio.TotalTarefas != 1 {
		t.Errorf("totais incorretos: %+v", relatorio)
	}

	if relatorio.PercentualMetasCumpridas != 100 {
		t.Errorf("esperava 100%% de metas cumpridas, obteve %v", relatorio.PercentualMetasCumpridas)
	}
}

func TestRelatorioHandlerGerarParametrosObrigatorios(t *testing.T) {
	handler := criarRelatorioHandlerTeste(
		&fakeRelatorioMetaRepository{},
		&fakeRelatorioTarefaRepository{},
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/relatorios",
		nil,
	)

	response := httptest.NewRecorder()

	handler.Gerar(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("esperava status 400, obteve %d", response.Code)
	}
}

func TestRelatorioHandlerGerarPeriodoInvalido(t *testing.T) {
	handler := criarRelatorioHandlerTeste(
		&fakeRelatorioMetaRepository{},
		&fakeRelatorioTarefaRepository{},
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/relatorios?data_inicio=2026-09-30&data_fim=2026-09-01",
		nil,
	)

	response := httptest.NewRecorder()

	handler.Gerar(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("esperava status 400, obteve %d", response.Code)
	}
}

func TestRelatorioHandlerGerarErroInterno(t *testing.T) {
	handler := criarRelatorioHandlerTeste(
		&fakeRelatorioMetaRepository{erroBuscarTodos: context.DeadlineExceeded},
		&fakeRelatorioTarefaRepository{},
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/relatorios?data_inicio=2026-09-01&data_fim=2026-09-30",
		nil,
	)

	response := httptest.NewRecorder()

	handler.Gerar(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("esperava status 500, obteve %d", response.Code)
	}
}
