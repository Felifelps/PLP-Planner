package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"plp-planner/models"
	"plp-planner/repositories"
	"plp-planner/services"
)

type fakeMetaRepository struct {
	meta  *models.Meta
	metas []models.Meta

	erroSalvar          error
	erroBuscarTodos     error
	erroBuscarPorID     error
	erroAtualizar       error
	erroAtualizarStatus error
	erroExcluir         error
}

func (f *fakeMetaRepository) Salvar(
	ctx context.Context,
	meta *models.Meta,
) error {
	if f.erroSalvar != nil {
		return f.erroSalvar
	}

	meta.ID = 1
	f.meta = meta

	return nil
}

func (f *fakeMetaRepository) BuscarTodos(
	ctx context.Context,
	dataInicio string,
	dataFim string,
) ([]models.Meta, error) {
	if f.erroBuscarTodos != nil {
		return nil, f.erroBuscarTodos
	}

	return f.metas, nil
}

func (f *fakeMetaRepository) BuscarPorID(
	ctx context.Context,
	id int64,
) (*models.Meta, error) {
	if f.erroBuscarPorID != nil {
		return nil, f.erroBuscarPorID
	}

	if f.meta == nil {
		return nil, repositories.ErrMetaNaoEncontrada
	}

	return f.meta, nil
}

func (f *fakeMetaRepository) Atualizar(
	ctx context.Context,
	meta *models.Meta,
) error {
	if f.erroAtualizar != nil {
		return f.erroAtualizar
	}

	f.meta = meta

	return nil
}

func (f *fakeMetaRepository) AtualizarStatus(
	ctx context.Context,
	id int64,
	status models.Status,
) error {
	if f.erroAtualizarStatus != nil {
		return f.erroAtualizarStatus
	}

	if f.meta == nil {
		return repositories.ErrMetaNaoEncontrada
	}

	f.meta.Status = status

	return nil
}

func (f *fakeMetaRepository) Excluir(
	ctx context.Context,
	id int64,
) error {
	if f.erroExcluir != nil {
		return f.erroExcluir
	}

	if f.meta == nil {
		return repositories.ErrMetaNaoEncontrada
	}

	f.meta = nil

	return nil
}

func criarMetaTeste() *models.Meta {
	return &models.Meta{
		ID:          1,
		Nome:        "Estudar algoritmos",
		Descricao:   "Estudar grafos",
		CategoriaID: 1,
		Status:      models.StatusCumprida,
		Periodo:     models.PeriodoSemanal,
		DataInicio: time.Date(
			2026,
			time.August,
			18,
			0,
			0,
			0,
			0,
			time.UTC,
		),
		DataFim: time.Date(
			2026,
			time.August,
			25,
			0,
			0,
			0,
			0,
			time.UTC,
		),
	}
}

func criarHandlerTeste(
	repository *fakeMetaRepository,
) *MetaHandler {
	service := services.NewMetaService(repository)
	return NewMetaHandler(service)
}

func TestMetaHandlerCriar(t *testing.T) {
	repository := &fakeMetaRepository{}
	handler := criarHandlerTeste(repository)

	body := `{
		"nome": "Estudar algoritmos",
		"descricao": "Estudar grafos",
		"categoria_id": 1,
		"status": "cumprida",
		"periodo": "semanal",
		"data_inicio": "2026-08-18T00:00:00Z",
		"data_fim": "2026-08-25T00:00:00Z"
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/metas",
		strings.NewReader(body),
	)

	response := httptest.NewRecorder()

	handler.Criar(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf(
			"status = %d; esperado %d",
			response.Code,
			http.StatusCreated,
		)
	}

	if repository.meta == nil {
		t.Fatal("meta não foi criada")
	}

	if repository.meta.Periodo != models.PeriodoSemanal {
		t.Fatalf(
			"período = %q; esperado %q",
			repository.meta.Periodo,
			models.PeriodoSemanal,
		)
	}
}

func TestMetaHandlerCriarDadosInvalidos(t *testing.T) {
	repository := &fakeMetaRepository{}
	handler := criarHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodPost,
		"/metas",
		strings.NewReader(`{"nome":`),
	)

	response := httptest.NewRecorder()

	handler.Criar(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"status = %d; esperado %d",
			response.Code,
			http.StatusBadRequest,
		)
	}
}

func TestMetaHandlerBuscarTodos(t *testing.T) {
	repository := &fakeMetaRepository{
		metas: []models.Meta{
			*criarMetaTeste(),
		},
	}

	handler := criarHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/metas?data_inicio=2026-08-01&data_fim=2026-08-31",
		nil,
	)

	response := httptest.NewRecorder()

	handler.BuscarTodos(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"status = %d; esperado %d",
			response.Code,
			http.StatusOK,
		)
	}

	var metas []models.Meta

	if err := json.NewDecoder(response.Body).Decode(&metas); err != nil {
		t.Fatalf(
			"erro ao decodificar resposta: %v",
			err,
		)
	}

	if len(metas) != 1 {
		t.Fatalf(
			"quantidade = %d; esperada 1",
			len(metas),
		)
	}
}

func TestMetaHandlerBuscarPorID(t *testing.T) {
	repository := &fakeMetaRepository{
		meta: criarMetaTeste(),
	}

	handler := criarHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/metas/1",
		nil,
	)

	response := httptest.NewRecorder()

	handler.BuscarPorID(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"status = %d; esperado %d",
			response.Code,
			http.StatusOK,
		)
	}

	var meta models.Meta

	if err := json.NewDecoder(response.Body).Decode(&meta); err != nil {
		t.Fatalf(
			"erro ao decodificar resposta: %v",
			err,
		)
	}

	if meta.ID != 1 {
		t.Fatalf(
			"id = %d; esperado 1",
			meta.ID,
		)
	}
}

func TestMetaHandlerBuscarPorIDInvalido(t *testing.T) {
	repository := &fakeMetaRepository{}
	handler := criarHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/metas/abc",
		nil,
	)

	response := httptest.NewRecorder()

	handler.BuscarPorID(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"status = %d; esperado %d",
			response.Code,
			http.StatusBadRequest,
		)
	}
}

func TestMetaHandlerBuscarPorIDNaoEncontrada(t *testing.T) {
	repository := &fakeMetaRepository{
		erroBuscarPorID: repositories.ErrMetaNaoEncontrada,
	}

	handler := criarHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/metas/1",
		nil,
	)

	response := httptest.NewRecorder()

	handler.BuscarPorID(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf(
			"status = %d; esperado %d",
			response.Code,
			http.StatusNotFound,
		)
	}
}

func TestMetaHandlerAtualizar(t *testing.T) {
	repository := &fakeMetaRepository{}

	handler := criarHandlerTeste(repository)

	body := `{
		"nome": "Estudar programação",
		"descricao": "Estudar Go",
		"categoria_id": 2,
		"status": "parcialmente cumprida",
		"periodo": "mensal",
		"data_inicio": "2026-08-01T00:00:00Z",
		"data_fim": "2026-08-31T00:00:00Z"
	}`

	request := httptest.NewRequest(
		http.MethodPut,
		"/metas/1",
		strings.NewReader(body),
	)

	response := httptest.NewRecorder()

	handler.Atualizar(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"status = %d; esperado %d",
			response.Code,
			http.StatusOK,
		)
	}

	if repository.meta == nil {
		t.Fatal("meta não foi atualizada")
	}

	if repository.meta.ID != 1 {
		t.Fatalf(
			"id = %d; esperado 1",
			repository.meta.ID,
		)
	}

	if repository.meta.Periodo != models.PeriodoMensal {
		t.Fatalf(
			"período = %q; esperado %q",
			repository.meta.Periodo,
			models.PeriodoMensal,
		)
	}
}

func TestMetaHandlerAtualizarStatus(t *testing.T) {
	repository := &fakeMetaRepository{
		meta: criarMetaTeste(),
	}

	handler := criarHandlerTeste(repository)

	body := `{
		"status": "não cumprida"
	}`

	request := httptest.NewRequest(
		http.MethodPatch,
		"/metas/1/status",
		strings.NewReader(body),
	)

	response := httptest.NewRecorder()

	handler.AtualizarStatus(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf(
			"status = %d; esperado %d",
			response.Code,
			http.StatusNoContent,
		)
	}

	if repository.meta.Status != models.StatusNaoCumprida {
		t.Fatalf(
			"status = %q; esperado %q",
			repository.meta.Status,
			models.StatusNaoCumprida,
		)
	}
}

func TestMetaHandlerAtualizarStatusInvalido(t *testing.T) {
	repository := &fakeMetaRepository{
		meta: criarMetaTeste(),
	}

	handler := criarHandlerTeste(repository)

	body := `{
		"status": "status inválido"
	}`

	request := httptest.NewRequest(
		http.MethodPatch,
		"/metas/1/status",
		strings.NewReader(body),
	)

	response := httptest.NewRecorder()

	handler.AtualizarStatus(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"status = %d; esperado %d",
			response.Code,
			http.StatusBadRequest,
		)
	}
}

func TestMetaHandlerExcluir(t *testing.T) {
	repository := &fakeMetaRepository{
		meta: criarMetaTeste(),
	}

	handler := criarHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/metas/1",
		nil,
	)

	response := httptest.NewRecorder()

	handler.Excluir(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf(
			"status = %d; esperado %d",
			response.Code,
			http.StatusNoContent,
		)
	}

	if repository.meta != nil {
		t.Fatal("meta não foi excluída")
	}
}

func TestMetaHandlerExcluirNaoEncontrada(t *testing.T) {
	repository := &fakeMetaRepository{
		erroExcluir: repositories.ErrMetaNaoEncontrada,
	}

	handler := criarHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/metas/1",
		nil,
	)

	response := httptest.NewRecorder()

	handler.Excluir(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf(
			"status = %d; esperado %d",
			response.Code,
			http.StatusNotFound,
		)
	}
}

func TestMetaHandlerBuscarTodosErroInterno(t *testing.T) {
	repository := &fakeMetaRepository{
		erroBuscarTodos: errors.New("erro no banco"),
	}

	handler := criarHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/metas",
		nil,
	)

	response := httptest.NewRecorder()

	handler.BuscarTodos(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf(
			"status = %d; esperado %d",
			response.Code,
			http.StatusInternalServerError,
		)
	}
}
