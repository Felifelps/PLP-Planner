package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"plp-planner/models"
	"plp-planner/repositories"
	"plp-planner/services"
)

type fakeTarefaRepository struct {
	tarefa  *models.Tarefa
	tarefas []models.Tarefa

	erroSalvar          error
	erroBuscarTodos     error
	erroBuscarPorID     error
	erroAtualizar       error
	erroAtualizarStatus error
	erroExcluir         error
}

func (f *fakeTarefaRepository) Salvar(
	ctx context.Context,
	tarefa *models.Tarefa,
) error {
	if f.erroSalvar != nil {
		return f.erroSalvar
	}

	tarefa.ID = 1
	f.tarefa = tarefa

	return nil
}

func (f *fakeTarefaRepository) BuscarTodos(
	ctx context.Context,
	data string,
	categoriaID string,
) ([]models.Tarefa, error) {
	if f.erroBuscarTodos != nil {
		return nil, f.erroBuscarTodos
	}

	return f.tarefas, nil
}

func (f *fakeTarefaRepository) BuscarPorID(
	ctx context.Context,
	id int64,
) (*models.Tarefa, error) {
	if f.erroBuscarPorID != nil {
		return nil, f.erroBuscarPorID
	}

	if f.tarefa == nil {
		return nil, repositories.ErrTarefaNaoEncontrada
	}

	return f.tarefa, nil
}

func (f *fakeTarefaRepository) Atualizar(
	ctx context.Context,
	tarefa *models.Tarefa,
) error {
	if f.erroAtualizar != nil {
		return f.erroAtualizar
	}

	f.tarefa = tarefa

	return nil
}

func (f *fakeTarefaRepository) AtualizarStatus(
	ctx context.Context,
	id int64,
	status models.StatusTarefa,
) error {
	if f.erroAtualizarStatus != nil {
		return f.erroAtualizarStatus
	}

	if f.tarefa == nil {
		return repositories.ErrTarefaNaoEncontrada
	}

	f.tarefa.Status = status

	return nil
}

func (f *fakeTarefaRepository) Excluir(
	ctx context.Context,
	id int64,
) error {
	if f.erroExcluir != nil {
		return f.erroExcluir
	}

	if f.tarefa == nil {
		return repositories.ErrTarefaNaoEncontrada
	}

	f.tarefa = nil

	return nil
}

func criarTarefaTeste() *models.Tarefa {
	turno := models.TurnoManha

	return &models.Tarefa{
		ID:          1,
		Descricao:   "Estudar Go",
		CategoriaID: 1,
		Data: time.Date(
			2026,
			time.August,
			19,
			0,
			0,
			0,
			0,
			time.UTC,
		),
		Turno:      &turno,
		Status:     models.StatusExecutada,
		Prioridade: models.PrioridadeAlta,
	}
}

func criarTarefaHandlerTeste(
	repository *fakeTarefaRepository,
) *TarefaHandler {
	service := services.NewTarefaService(repository)
	return NewTarefaHandler(service)
}

func TestTarefaHandlerCriar(t *testing.T) {
	repository := &fakeTarefaRepository{}
	handler := criarTarefaHandlerTeste(repository)

	body := `{
		"descricao": "Estudar Go",
		"categoria_id": 1,
		"data": "2026-08-19T00:00:00Z",
		"turno": "manhã",
		"status": "executada",
		"prioridade": "alta"
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/tarefas",
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

	if repository.tarefa == nil {
		t.Fatal("tarefa não foi criada")
	}
}

func TestTarefaHandlerCriarDadosInvalidos(t *testing.T) {
	repository := &fakeTarefaRepository{}
	handler := criarTarefaHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodPost,
		"/tarefas",
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

func TestTarefaHandlerBuscarTodos(t *testing.T) {
	repository := &fakeTarefaRepository{
		tarefas: []models.Tarefa{
			*criarTarefaTeste(),
		},
	}

	handler := criarTarefaHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/tarefas?data=2026-08-19",
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

	var tarefas []models.Tarefa

	if err := json.NewDecoder(response.Body).Decode(&tarefas); err != nil {
		t.Fatalf(
			"erro ao decodificar resposta: %v",
			err,
		)
	}

	if len(tarefas) != 1 {
		t.Fatalf(
			"quantidade = %d; esperada 1",
			len(tarefas),
		)
	}
}

func TestTarefaHandlerBuscarPorID(t *testing.T) {
	repository := &fakeTarefaRepository{
		tarefa: criarTarefaTeste(),
	}

	handler := criarTarefaHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/tarefas/1",
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
}

func TestTarefaHandlerBuscarPorIDNaoEncontrada(t *testing.T) {
	repository := &fakeTarefaRepository{
		erroBuscarPorID: repositories.ErrTarefaNaoEncontrada,
	}

	handler := criarTarefaHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/tarefas/1",
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

func TestTarefaHandlerAtualizarStatus(t *testing.T) {
	repository := &fakeTarefaRepository{
		tarefa: criarTarefaTeste(),
	}

	handler := criarTarefaHandlerTeste(repository)

	body := `{
		"status": "adiada"
	}`

	request := httptest.NewRequest(
		http.MethodPatch,
		"/tarefas/1/status",
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

	if repository.tarefa.Status != models.StatusAdiada {
		t.Fatalf(
			"status = %q; esperado %q",
			repository.tarefa.Status,
			models.StatusAdiada,
		)
	}
}

func TestTarefaHandlerAtualizarStatusInvalido(t *testing.T) {
	repository := &fakeTarefaRepository{
		tarefa: criarTarefaTeste(),
	}

	handler := criarTarefaHandlerTeste(repository)

	body := `{
		"status": "status inválido"
	}`

	request := httptest.NewRequest(
		http.MethodPatch,
		"/tarefas/1/status",
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

func TestTarefaHandlerExcluir(t *testing.T) {
	repository := &fakeTarefaRepository{
		tarefa: criarTarefaTeste(),
	}

	handler := criarTarefaHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/tarefas/1",
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

	if repository.tarefa != nil {
		t.Fatal("tarefa não foi excluída")
	}
}

func TestTarefaHandlerExcluirNaoEncontrada(t *testing.T) {
	repository := &fakeTarefaRepository{
		erroExcluir: repositories.ErrTarefaNaoEncontrada,
	}

	handler := criarTarefaHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/tarefas/1",
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
