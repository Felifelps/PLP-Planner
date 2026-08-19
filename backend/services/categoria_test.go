package services

import (
	"context"
	"errors"
	"testing"

	"plp-planner/models"
)

type fakeCategoriaRepository struct {
	salvarChamado      bool
	buscarTodosChamado bool
	buscarPorIDChamado bool
	atualizarChamado   bool
	excluirChamado     bool

	categoriaRecebida *models.Categoria
	idRecebido        int64

	categoriasRetornadas []models.Categoria
	categoriaRetornada   *models.Categoria

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
	f.salvarChamado = true
	f.categoriaRecebida = categoria
	return f.erroSalvar
}

func (f *fakeCategoriaRepository) BuscarTodos(
	ctx context.Context,
) ([]models.Categoria, error) {
	f.buscarTodosChamado = true
	return f.categoriasRetornadas, f.erroBuscarTodos
}

func (f *fakeCategoriaRepository) BuscarPorID(
	ctx context.Context,
	id int64,
) (*models.Categoria, error) {
	f.buscarPorIDChamado = true
	f.idRecebido = id
	return f.categoriaRetornada, f.erroBuscarPorID
}

func (f *fakeCategoriaRepository) Atualizar(
	ctx context.Context,
	categoria *models.Categoria,
) error {
	f.atualizarChamado = true
	f.categoriaRecebida = categoria
	return f.erroAtualizar
}

func (f *fakeCategoriaRepository) Excluir(
	ctx context.Context,
	id int64,
) error {
	f.excluirChamado = true
	f.idRecebido = id
	return f.erroExcluir
}

func criarCategoriaValida() *models.Categoria {
	return &models.Categoria{
		ID:   1,
		Nome: "Trabalho",
		Cor:  "#4C6EF5",
	}
}

func TestCategoriaServiceSalvar(t *testing.T) {
	repository := &fakeCategoriaRepository{}
	service := NewCategoriaService(repository)

	categoria := criarCategoriaValida()

	err := service.Salvar(
		context.Background(),
		categoria,
	)

	if err != nil {
		t.Fatalf(
			"Salvar() retornou erro inesperado: %v",
			err,
		)
	}

	if !repository.salvarChamado {
		t.Fatal("Salvar() do repository não foi chamado")
	}

	if repository.categoriaRecebida != categoria {
		t.Fatal("a categoria recebida pelo repository é diferente da enviada")
	}
}

func TestCategoriaServiceSalvarCategoriaNula(t *testing.T) {
	repository := &fakeCategoriaRepository{}
	service := NewCategoriaService(repository)

	err := service.Salvar(
		context.Background(),
		nil,
	)

	if !errors.Is(err, ErrCategoriaInvalidaEntrada) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrCategoriaInvalidaEntrada,
		)
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestCategoriaServiceSalvarNomeVazio(t *testing.T) {
	repository := &fakeCategoriaRepository{}
	service := NewCategoriaService(repository)

	categoria := criarCategoriaValida()
	categoria.Nome = "   "

	err := service.Salvar(
		context.Background(),
		categoria,
	)

	if !errors.Is(err, ErrCategoriaNomeObrigatorio) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrCategoriaNomeObrigatorio,
		)
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestCategoriaServiceSalvarCorInvalida(t *testing.T) {
	repository := &fakeCategoriaRepository{}
	service := NewCategoriaService(repository)

	categoria := criarCategoriaValida()
	categoria.Cor = "invalida"

	err := service.Salvar(
		context.Background(),
		categoria,
	)

	if err == nil {
		t.Fatal("esperava erro para cor inválida")
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestCategoriaServiceBuscarPorID(t *testing.T) {
	categoriaEsperada := criarCategoriaValida()

	repository := &fakeCategoriaRepository{
		categoriaRetornada: categoriaEsperada,
	}

	service := NewCategoriaService(repository)

	categoria, err := service.BuscarPorID(
		context.Background(),
		1,
	)

	if err != nil {
		t.Fatalf(
			"BuscarPorID() retornou erro inesperado: %v",
			err,
		)
	}

	if !repository.buscarPorIDChamado {
		t.Fatal("BuscarPorID() do repository não foi chamado")
	}

	if categoria != categoriaEsperada {
		t.Fatal("categoria retornada é diferente da esperada")
	}
}

func TestCategoriaServiceBuscarPorIDIDInvalido(t *testing.T) {
	repository := &fakeCategoriaRepository{}
	service := NewCategoriaService(repository)

	_, err := service.BuscarPorID(
		context.Background(),
		0,
	)

	if !errors.Is(err, ErrCategoriaIDInvalido) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrCategoriaIDInvalido,
		)
	}

	if repository.buscarPorIDChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestCategoriaServiceBuscarTodos(t *testing.T) {
	categoriasEsperadas := []models.Categoria{
		*criarCategoriaValida(),
	}

	repository := &fakeCategoriaRepository{
		categoriasRetornadas: categoriasEsperadas,
	}

	service := NewCategoriaService(repository)

	categorias, err := service.BuscarTodos(context.Background())

	if err != nil {
		t.Fatalf(
			"BuscarTodos() retornou erro inesperado: %v",
			err,
		)
	}

	if !repository.buscarTodosChamado {
		t.Fatal("BuscarTodos() do repository não foi chamado")
	}

	if len(categorias) != 1 {
		t.Fatalf(
			"quantidade de categorias = %d; esperada 1",
			len(categorias),
		)
	}
}

func TestCategoriaServiceAtualizar(t *testing.T) {
	repository := &fakeCategoriaRepository{}
	service := NewCategoriaService(repository)

	categoria := criarCategoriaValida()

	err := service.Atualizar(
		context.Background(),
		categoria,
	)

	if err != nil {
		t.Fatalf(
			"Atualizar() retornou erro inesperado: %v",
			err,
		)
	}

	if !repository.atualizarChamado {
		t.Fatal("Atualizar() do repository não foi chamado")
	}
}

func TestCategoriaServiceAtualizarIDInvalido(t *testing.T) {
	repository := &fakeCategoriaRepository{}
	service := NewCategoriaService(repository)

	categoria := criarCategoriaValida()
	categoria.ID = 0

	err := service.Atualizar(
		context.Background(),
		categoria,
	)

	if !errors.Is(err, ErrCategoriaIDInvalido) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrCategoriaIDInvalido,
		)
	}

	if repository.atualizarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestCategoriaServiceExcluir(t *testing.T) {
	repository := &fakeCategoriaRepository{}
	service := NewCategoriaService(repository)

	err := service.Excluir(
		context.Background(),
		1,
	)

	if err != nil {
		t.Fatalf(
			"Excluir() retornou erro inesperado: %v",
			err,
		)
	}

	if !repository.excluirChamado {
		t.Fatal("Excluir() do repository não foi chamado")
	}
}

func TestCategoriaServiceExcluirIDInvalido(t *testing.T) {
	repository := &fakeCategoriaRepository{}
	service := NewCategoriaService(repository)

	err := service.Excluir(
		context.Background(),
		0,
	)

	if !errors.Is(err, ErrCategoriaIDInvalido) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrCategoriaIDInvalido,
		)
	}

	if repository.excluirChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}
