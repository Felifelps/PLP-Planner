package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"plp-planner/models"
)

type fakeMetaRepository struct {
	salvarChamado          bool
	buscarTodosChamado     bool
	buscarPorIDChamado     bool
	atualizarChamado       bool
	atualizarStatusChamado bool
	excluirChamado         bool

	metaRecebida   *models.Meta
	idRecebido     int64
	statusRecebido models.Status

	metasRetornadas []models.Meta
	metaRetornada   *models.Meta

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
	f.salvarChamado = true
	f.metaRecebida = meta
	return f.erroSalvar
}

func (f *fakeMetaRepository) BuscarTodos(
	ctx context.Context,
	dataInicio string,
	dataFim string,
) ([]models.Meta, error) {
	f.buscarTodosChamado = true
	return f.metasRetornadas, f.erroBuscarTodos
}

func (f *fakeMetaRepository) BuscarPorID(
	ctx context.Context,
	id int64,
) (*models.Meta, error) {
	f.buscarPorIDChamado = true
	f.idRecebido = id
	return f.metaRetornada, f.erroBuscarPorID
}

func (f *fakeMetaRepository) Atualizar(
	ctx context.Context,
	meta *models.Meta,
) error {
	f.atualizarChamado = true
	f.metaRecebida = meta
	return f.erroAtualizar
}

func (f *fakeMetaRepository) AtualizarStatus(
	ctx context.Context,
	id int64,
	status models.Status,
) error {
	f.atualizarStatusChamado = true
	f.idRecebido = id
	f.statusRecebido = status
	return f.erroAtualizarStatus
}

func (f *fakeMetaRepository) Excluir(
	ctx context.Context,
	id int64,
) error {
	f.excluirChamado = true
	f.idRecebido = id
	return f.erroExcluir
}

func criarMetaValida() *models.Meta {
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

func TestMetaServiceSalvar(t *testing.T) {
	repository := &fakeMetaRepository{}
	service := NewMetaService(repository)

	meta := criarMetaValida()

	err := service.Salvar(
		context.Background(),
		meta,
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

	if repository.metaRecebida != meta {
		t.Fatal("a meta recebida pelo repository é diferente da enviada")
	}
}

func TestMetaServiceSalvarMetaNula(t *testing.T) {
	repository := &fakeMetaRepository{}
	service := NewMetaService(repository)

	err := service.Salvar(
		context.Background(),
		nil,
	)

	if !errors.Is(err, ErrMetaInvalida) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrMetaInvalida,
		)
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestMetaServiceSalvarNomeVazio(t *testing.T) {
	repository := &fakeMetaRepository{}
	service := NewMetaService(repository)

	meta := criarMetaValida()
	meta.Nome = "   "

	err := service.Salvar(
		context.Background(),
		meta,
	)

	if !errors.Is(err, ErrNomeObrigatorio) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrNomeObrigatorio,
		)
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestMetaServiceSalvarCategoriaInvalida(t *testing.T) {
	repository := &fakeMetaRepository{}
	service := NewMetaService(repository)

	meta := criarMetaValida()
	meta.CategoriaID = 0

	err := service.Salvar(
		context.Background(),
		meta,
	)

	if !errors.Is(err, ErrCategoriaInvalida) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrCategoriaInvalida,
		)
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestMetaServiceSalvarStatusInvalido(t *testing.T) {
	repository := &fakeMetaRepository{}
	service := NewMetaService(repository)

	meta := criarMetaValida()
	meta.Status = models.Status("invalido")

	err := service.Salvar(
		context.Background(),
		meta,
	)

	if err == nil {
		t.Fatal("esperava erro para status inválido")
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestMetaServiceSalvarPeriodoInvalido(t *testing.T) {
	repository := &fakeMetaRepository{}
	service := NewMetaService(repository)

	meta := criarMetaValida()
	meta.Periodo = models.Periodo("invalido")

	err := service.Salvar(
		context.Background(),
		meta,
	)

	if err == nil {
		t.Fatal("esperava erro para período inválido")
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestMetaServiceSalvarDatasInvalidas(t *testing.T) {
	repository := &fakeMetaRepository{}
	service := NewMetaService(repository)

	meta := criarMetaValida()
	meta.DataInicio = meta.DataFim.AddDate(0, 0, 1)

	err := service.Salvar(
		context.Background(),
		meta,
	)

	if err == nil {
		t.Fatal("esperava erro para datas inválidas")
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestMetaServiceBuscarPorID(t *testing.T) {
	metaEsperada := criarMetaValida()

	repository := &fakeMetaRepository{
		metaRetornada: metaEsperada,
	}

	service := NewMetaService(repository)

	meta, err := service.BuscarPorID(
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

	if repository.idRecebido != 1 {
		t.Fatalf(
			"id recebido = %d; esperado 1",
			repository.idRecebido,
		)
	}

	if meta != metaEsperada {
		t.Fatal("meta retornada é diferente da esperada")
	}
}

func TestMetaServiceBuscarPorIDIDInvalido(t *testing.T) {
	repository := &fakeMetaRepository{}
	service := NewMetaService(repository)

	_, err := service.BuscarPorID(
		context.Background(),
		0,
	)

	if !errors.Is(err, ErrIDInvalido) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrIDInvalido,
		)
	}

	if repository.buscarPorIDChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestMetaServiceBuscarTodos(t *testing.T) {
	metasEsperadas := []models.Meta{
		*criarMetaValida(),
	}

	repository := &fakeMetaRepository{
		metasRetornadas: metasEsperadas,
	}

	service := NewMetaService(repository)

	metas, err := service.BuscarTodos(
		context.Background(),
		"2026-08-01",
		"2026-08-31",
	)

	if err != nil {
		t.Fatalf(
			"BuscarTodos() retornou erro inesperado: %v",
			err,
		)
	}

	if !repository.buscarTodosChamado {
		t.Fatal("BuscarTodos() do repository não foi chamado")
	}

	if len(metas) != 1 {
		t.Fatalf(
			"quantidade de metas = %d; esperada 1",
			len(metas),
		)
	}
}

func TestMetaServiceAtualizar(t *testing.T) {
	repository := &fakeMetaRepository{}
	service := NewMetaService(repository)

	meta := criarMetaValida()

	err := service.Atualizar(
		context.Background(),
		meta,
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

	if repository.metaRecebida != meta {
		t.Fatal("meta recebida pelo repository é diferente")
	}
}

func TestMetaServiceAtualizarMetaNula(t *testing.T) {
	repository := &fakeMetaRepository{}
	service := NewMetaService(repository)

	err := service.Atualizar(
		context.Background(),
		nil,
	)

	if !errors.Is(err, ErrMetaInvalida) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrMetaInvalida,
		)
	}

	if repository.atualizarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestMetaServiceAtualizarIDInvalido(t *testing.T) {
	repository := &fakeMetaRepository{}
	service := NewMetaService(repository)

	meta := criarMetaValida()
	meta.ID = 0

	err := service.Atualizar(
		context.Background(),
		meta,
	)

	if !errors.Is(err, ErrIDInvalido) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrIDInvalido,
		)
	}

	if repository.atualizarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestMetaServiceAtualizarStatus(t *testing.T) {
	repository := &fakeMetaRepository{}
	service := NewMetaService(repository)

	err := service.AtualizarStatus(
		context.Background(),
		1,
		models.StatusParcialmenteCumprida,
	)

	if err != nil {
		t.Fatalf(
			"AtualizarStatus() retornou erro inesperado: %v",
			err,
		)
	}

	if !repository.atualizarStatusChamado {
		t.Fatal("AtualizarStatus() do repository não foi chamado")
	}

	if repository.idRecebido != 1 {
		t.Fatalf(
			"id recebido = %d; esperado 1",
			repository.idRecebido,
		)
	}

	if repository.statusRecebido != models.StatusParcialmenteCumprida {
		t.Fatalf(
			"status recebido = %q; esperado %q",
			repository.statusRecebido,
			models.StatusParcialmenteCumprida,
		)
	}
}

func TestMetaServiceAtualizarStatusIDInvalido(t *testing.T) {
	repository := &fakeMetaRepository{}
	service := NewMetaService(repository)

	err := service.AtualizarStatus(
		context.Background(),
		0,
		models.StatusCumprida,
	)

	if !errors.Is(err, ErrIDInvalido) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrIDInvalido,
		)
	}

	if repository.atualizarStatusChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestMetaServiceAtualizarStatusInvalido(t *testing.T) {
	repository := &fakeMetaRepository{}
	service := NewMetaService(repository)

	err := service.AtualizarStatus(
		context.Background(),
		1,
		models.Status("invalido"),
	)

	if err == nil {
		t.Fatal("esperava erro para status inválido")
	}

	if repository.atualizarStatusChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestMetaServiceExcluir(t *testing.T) {
	repository := &fakeMetaRepository{}
	service := NewMetaService(repository)

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

	if repository.idRecebido != 1 {
		t.Fatalf(
			"id recebido = %d; esperado 1",
			repository.idRecebido,
		)
	}
}

func TestMetaServiceExcluirIDInvalido(t *testing.T) {
	repository := &fakeMetaRepository{}
	service := NewMetaService(repository)

	err := service.Excluir(
		context.Background(),
		0,
	)

	if !errors.Is(err, ErrIDInvalido) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrIDInvalido,
		)
	}

	if repository.excluirChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}
