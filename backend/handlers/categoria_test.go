package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"plp-planner/models"
	"plp-planner/repositories"
	"plp-planner/services"
)

type fakeCategoriaRepository struct {
	categoria  *models.Categoria
	categorias []models.Categoria

	erroSalvar      error
	erroBuscarTodos error
	erroBuscarPorID error
	erroAtualizar   error
	erroExcluir     error
}

func (f *fakeCategoriaRepository) Salvar(
	ctx context.Context,
	categoria *models.Categoria,
) error {
	if f.erroSalvar != nil {
		return f.erroSalvar
	}

	categoria.ID = 1
	f.categoria = categoria

	return nil
}

func (f *fakeCategoriaRepository) BuscarTodos(
	ctx context.Context,
) ([]models.Categoria, error) {
	if f.erroBuscarTodos != nil {
		return nil, f.erroBuscarTodos
	}

	return f.categorias, nil
}

func (f *fakeCategoriaRepository) BuscarPorID(
	ctx context.Context,
	id int64,
) (*models.Categoria, error) {
	if f.erroBuscarPorID != nil {
		return nil, f.erroBuscarPorID
	}

	if f.categoria == nil {
		return nil, repositories.ErrCategoriaNaoEncontrada
	}

	return f.categoria, nil
}

func (f *fakeCategoriaRepository) Atualizar(
	ctx context.Context,
	categoria *models.Categoria,
) error {
	if f.erroAtualizar != nil {
		return f.erroAtualizar
	}

	f.categoria = categoria

	return nil
}

func (f *fakeCategoriaRepository) Excluir(
	ctx context.Context,
	id int64,
) error {
	if f.erroExcluir != nil {
		return f.erroExcluir
	}

	if f.categoria == nil {
		return repositories.ErrCategoriaNaoEncontrada
	}

	f.categoria = nil

	return nil
}

func criarCategoriaTeste() *models.Categoria {
	return &models.Categoria{
		ID:   1,
		Nome: "Trabalho",
		Cor:  "#4C6EF5",
	}
}

func criarCategoriaHandlerTeste(
	repository *fakeCategoriaRepository,
) *CategoriaHandler {
	service := services.NewCategoriaService(repository)
	return NewCategoriaHandler(service)
}

func TestCategoriaHandlerCriar(t *testing.T) {
	repository := &fakeCategoriaRepository{}
	handler := criarCategoriaHandlerTeste(repository)

	body := `{
		"nome": "Trabalho",
		"cor": "#4C6EF5"
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/categorias",
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

	if repository.categoria == nil {
		t.Fatal("categoria não foi criada")
	}
}

func TestCategoriaHandlerCriarDadosInvalidos(t *testing.T) {
	repository := &fakeCategoriaRepository{}
	handler := criarCategoriaHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodPost,
		"/categorias",
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

func TestCategoriaHandlerBuscarTodos(t *testing.T) {
	repository := &fakeCategoriaRepository{
		categorias: []models.Categoria{
			*criarCategoriaTeste(),
		},
	}

	handler := criarCategoriaHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/categorias",
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

	var categorias []models.Categoria

	if err := json.NewDecoder(response.Body).Decode(&categorias); err != nil {
		t.Fatalf(
			"erro ao decodificar resposta: %v",
			err,
		)
	}

	if len(categorias) != 1 {
		t.Fatalf(
			"quantidade = %d; esperada 1",
			len(categorias),
		)
	}
}

func TestCategoriaHandlerBuscarPorID(t *testing.T) {
	repository := &fakeCategoriaRepository{
		categoria: criarCategoriaTeste(),
	}

	handler := criarCategoriaHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/categorias/1",
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

func TestCategoriaHandlerBuscarPorIDInvalido(t *testing.T) {
	repository := &fakeCategoriaRepository{}
	handler := criarCategoriaHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/categorias/abc",
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

func TestCategoriaHandlerBuscarPorIDNaoEncontrada(t *testing.T) {
	repository := &fakeCategoriaRepository{
		erroBuscarPorID: repositories.ErrCategoriaNaoEncontrada,
	}

	handler := criarCategoriaHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/categorias/1",
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

func TestCategoriaHandlerAtualizar(t *testing.T) {
	repository := &fakeCategoriaRepository{}

	handler := criarCategoriaHandlerTeste(repository)

	body := `{
		"nome": "Estudos",
		"cor": "#7048E8"
	}`

	request := httptest.NewRequest(
		http.MethodPut,
		"/categorias/1",
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

	if repository.categoria == nil {
		t.Fatal("categoria não foi atualizada")
	}
}

func TestCategoriaHandlerExcluir(t *testing.T) {
	repository := &fakeCategoriaRepository{
		categoria: criarCategoriaTeste(),
	}

	handler := criarCategoriaHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/categorias/1",
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

	if repository.categoria != nil {
		t.Fatal("categoria não foi excluída")
	}
}

func TestCategoriaHandlerExcluirNaoEncontrada(t *testing.T) {
	repository := &fakeCategoriaRepository{
		erroExcluir: repositories.ErrCategoriaNaoEncontrada,
	}

	handler := criarCategoriaHandlerTeste(repository)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/categorias/1",
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
