package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"plp-planner/models"
	"plp-planner/repositories"
	"plp-planner/services"
)

type fakeLembreteRepository struct {
	lembrete  *models.Lembrete
	lembretes []models.Lembrete

	erroSalvar      error
	erroBuscarTodos error
	erroBuscarPorID error
	erroAtualizar   error
	erroExcluir     error
}

func (f *fakeLembreteRepository) Salvar(
	ctx context.Context,
	lembrete *models.Lembrete,
) error {
	if f.erroSalvar != nil {
		return f.erroSalvar
	}

	lembrete.ID = 1
	f.lembrete = lembrete

	return nil
}

func (f *fakeLembreteRepository) BuscarTodos(
	ctx context.Context,
	dataInicio string,
	dataFim string,
) ([]models.Lembrete, error) {
	if f.erroBuscarTodos != nil {
		return nil, f.erroBuscarTodos
	}

	return f.lembretes, nil
}

func (f *fakeLembreteRepository) BuscarPorID(
	ctx context.Context,
	id int64,
) (*models.Lembrete, error) {
	if f.erroBuscarPorID != nil {
		return nil, f.erroBuscarPorID
	}

	if f.lembrete == nil {
		return nil, repositories.ErrLembreteNaoEncontrado
	}

	return f.lembrete, nil
}

func (f *fakeLembreteRepository) Atualizar(
	ctx context.Context,
	lembrete *models.Lembrete,
) error {
	if f.erroAtualizar != nil {
		return f.erroAtualizar
	}

	f.lembrete = lembrete

	return nil
}

func (f *fakeLembreteRepository) Excluir(
	ctx context.Context,
	id int64,
) error {
	if f.erroExcluir != nil {
		return f.erroExcluir
	}

	if f.lembrete == nil {
		return repositories.ErrLembreteNaoEncontrado
	}

	f.lembrete = nil

	return nil
}

func criarLembreteTeste() *models.Lembrete {
	return &models.Lembrete{
		ID:         1,
		Descricao:  "Entregar trabalho de PLP",
		Tipo:       models.TipoEntrega,
		Data:       "2026-08-25",
		Horario:    "14:30",
		Recorrente: false,
	}
}

func criarLembreteHandlerTeste(
	repository *fakeLembreteRepository,
) *LembreteHandler {
	service := services.NewLembreteService(repository)

	return NewLembreteHandler(service)
}

func TestLembreteHandlerCriar(t *testing.T) {
	repository := &fakeLembreteRepository{}
	handler := criarLembreteHandlerTeste(repository)

	body := `{
		"descricao": "Entregar trabalho de PLP",
		"tipo": "entrega",
		"data": "2026-08-25",
		"horario": "14:30",
		"recorrente": false
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/lembretes",
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

	if repository.lembrete == nil {
		t.Fatal("lembrete não foi criado")
	}

	if repository.lembrete.ID != 1 {
		t.Fatalf(
			"id = %d; esperado 1",
			repository.lembrete.ID,
		)
	}

	if repository.lembrete.Tipo != models.TipoEntrega {
		t.Fatalf(
			"tipo = %q; esperado %q",
			repository.lembrete.Tipo,
			models.TipoEntrega,
		)
	}

	if repository.lembrete.Recorrente {
		t.Fatal("lembrete não deveria ser recorrente")
	}
}

func TestLembreteHandlerCriarRecorrente(t *testing.T) {
	repository := &fakeLembreteRepository{}
	handler := criarLembreteHandlerTeste(repository)

	body := `{
		"descricao": "Estudar PLP",
		"tipo": "estudo",
		"data": "2026-08-24",
		"horario": "19:00",
		"recorrente": true
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/lembretes",
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

	if repository.lembrete == nil {
		t.Fatal("lembrete não foi criado")
	}

	if !repository.lembrete.Recorrente {
		t.Fatal("lembrete deveria ser recorrente")
	}
}

func TestLembreteHandlerCriarDadosInvalidos(t *testing.T) {
	repository := &fakeLembreteRepository{}
	handler := criarLembreteHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodPost,
		"/lembretes",
		strings.NewReader(`{"descricao":`),
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

func TestLembreteHandlerCriarDescricaoVazia(t *testing.T) {
	repository := &fakeLembreteRepository{}
	handler := criarLembreteHandlerTeste(repository)

	body := `{
		"descricao": "",
		"tipo": "estudo",
		"data": "2026-08-25",
		"horario": "14:30",
		"recorrente": false
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/lembretes",
		strings.NewReader(body),
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

func TestLembreteHandlerBuscarTodos(t *testing.T) {
	repository := &fakeLembreteRepository{
		lembretes: []models.Lembrete{
			*criarLembreteTeste(),
		},
	}

	handler := criarLembreteHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/lembretes?data_inicio=2026-08-01&data_fim=2026-08-31",
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

	var lembretes []models.Lembrete

	if err := json.NewDecoder(response.Body).Decode(
		&lembretes,
	); err != nil {
		t.Fatalf(
			"erro ao decodificar resposta: %v",
			err,
		)
	}

	if len(lembretes) != 1 {
		t.Fatalf(
			"quantidade = %d; esperada 1",
			len(lembretes),
		)
	}

	if lembretes[0].Descricao != "Entregar trabalho de PLP" {
		t.Fatalf(
			"descrição = %q; esperada %q",
			lembretes[0].Descricao,
			"Entregar trabalho de PLP",
		)
	}
}

func TestLembreteHandlerBuscarTodosRecorrente(t *testing.T) {
	repository := &fakeLembreteRepository{
		lembretes: []models.Lembrete{
			{
				ID:         1,
				Descricao:  "Estudar PLP",
				Tipo:       models.TipoEstudo,
				Data:       "2026-08-03",
				Horario:    "19:00",
				Recorrente: true,
			},
		},
	}

	handler := criarLembreteHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/lembretes?data_inicio=2026-08-17&data_fim=2026-08-31",
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

	var lembretes []models.Lembrete

	if err := json.NewDecoder(response.Body).Decode(
		&lembretes,
	); err != nil {
		t.Fatalf(
			"erro ao decodificar resposta: %v",
			err,
		)
	}

	if len(lembretes) != 3 {
		t.Fatalf(
			"quantidade = %d; esperada 3",
			len(lembretes),
		)
	}

	datasEsperadas := []string{
		"2026-08-17",
		"2026-08-24",
		"2026-08-31",
	}

	for i, data := range datasEsperadas {
		if lembretes[i].Data != data {
			t.Fatalf(
				"data = %s; esperada %s",
				lembretes[i].Data,
				data,
			)
		}
	}
}

func TestLembreteHandlerBuscarTodosPeriodoInvalido(t *testing.T) {
	repository := &fakeLembreteRepository{}
	handler := criarLembreteHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/lembretes?data_inicio=2026-08-31&data_fim=2026-08-01",
		nil,
	)

	response := httptest.NewRecorder()

	handler.BuscarTodos(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"status = %d; esperado %d",
			response.Code,
			http.StatusBadRequest,
		)
	}
}

func TestLembreteHandlerBuscarPorID(t *testing.T) {
	repository := &fakeLembreteRepository{
		lembrete: criarLembreteTeste(),
	}

	handler := criarLembreteHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/lembretes/1",
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

	var lembrete models.Lembrete

	if err := json.NewDecoder(response.Body).Decode(
		&lembrete,
	); err != nil {
		t.Fatalf(
			"erro ao decodificar resposta: %v",
			err,
		)
	}

	if lembrete.ID != 1 {
		t.Fatalf(
			"id = %d; esperado 1",
			lembrete.ID,
		)
	}
}

func TestLembreteHandlerBuscarPorIDInvalido(t *testing.T) {
	repository := &fakeLembreteRepository{}
	handler := criarLembreteHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/lembretes/abc",
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

func TestLembreteHandlerBuscarPorIDNaoEncontrado(t *testing.T) {
	repository := &fakeLembreteRepository{
		erroBuscarPorID: repositories.ErrLembreteNaoEncontrado,
	}

	handler := criarLembreteHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/lembretes/1",
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

func TestLembreteHandlerAtualizar(t *testing.T) {
	repository := &fakeLembreteRepository{}
	handler := criarLembreteHandlerTeste(repository)

	body := `{
		"descricao": "Reunião do projeto",
		"tipo": "reunião",
		"data": "2026-08-28",
		"horario": "10:00",
		"recorrente": true
	}`

	request := httptest.NewRequest(
		http.MethodPut,
		"/lembretes/1",
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

	if repository.lembrete == nil {
		t.Fatal("lembrete não foi atualizado")
	}

	if repository.lembrete.ID != 1 {
		t.Fatalf(
			"id = %d; esperado 1",
			repository.lembrete.ID,
		)
	}

	if repository.lembrete.Tipo != models.TipoReuniao {
		t.Fatalf(
			"tipo = %q; esperado %q",
			repository.lembrete.Tipo,
			models.TipoReuniao,
		)
	}

	if !repository.lembrete.Recorrente {
		t.Fatal("lembrete deveria ser recorrente")
	}
}

func TestLembreteHandlerAtualizarNaoEncontrado(t *testing.T) {
	repository := &fakeLembreteRepository{
		erroAtualizar: repositories.ErrLembreteNaoEncontrado,
	}

	handler := criarLembreteHandlerTeste(repository)

	body := `{
		"descricao": "Reunião",
		"tipo": "reunião",
		"data": "2026-08-28",
		"horario": "10:00",
		"recorrente": false
	}`

	request := httptest.NewRequest(
		http.MethodPut,
		"/lembretes/1",
		strings.NewReader(body),
	)

	response := httptest.NewRecorder()

	handler.Atualizar(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf(
			"status = %d; esperado %d",
			response.Code,
			http.StatusNotFound,
		)
	}
}

func TestLembreteHandlerExcluir(t *testing.T) {
	repository := &fakeLembreteRepository{
		lembrete: criarLembreteTeste(),
	}

	handler := criarLembreteHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/lembretes/1",
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

	if repository.lembrete != nil {
		t.Fatal("lembrete não foi excluído")
	}
}

func TestLembreteHandlerExcluirNaoEncontrado(t *testing.T) {
	repository := &fakeLembreteRepository{
		erroExcluir: repositories.ErrLembreteNaoEncontrado,
	}

	handler := criarLembreteHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/lembretes/1",
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

func TestLembreteHandlerBuscarTodosErroInterno(t *testing.T) {
	repository := &fakeLembreteRepository{
		erroBuscarTodos: errors.New("erro no banco"),
	}

	handler := criarLembreteHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/lembretes",
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