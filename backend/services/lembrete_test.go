package services

import (
	"context"
	"errors"
	"testing"

	"plp-planner/models"
)

type fakeLembreteRepository struct {
	salvarChamado      bool
	buscarTodosChamado bool
	buscarPorIDChamado bool
	atualizarChamado   bool
	excluirChamado     bool

	lembreteRecebido *models.Lembrete
	idRecebido       int64

	dataInicioRecebida string
	dataFimRecebida    string

	lembretesRetornados []models.Lembrete
	lembreteRetornado   *models.Lembrete

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
	f.salvarChamado = true
	f.lembreteRecebido = lembrete
	return f.erroSalvar
}

func (f *fakeLembreteRepository) BuscarTodos(
	ctx context.Context,
	dataInicio string,
	dataFim string,
) ([]models.Lembrete, error) {
	f.buscarTodosChamado = true
	f.dataInicioRecebida = dataInicio
	f.dataFimRecebida = dataFim

	return f.lembretesRetornados, f.erroBuscarTodos
}

func (f *fakeLembreteRepository) BuscarPorID(
	ctx context.Context,
	id int64,
) (*models.Lembrete, error) {
	f.buscarPorIDChamado = true
	f.idRecebido = id

	return f.lembreteRetornado, f.erroBuscarPorID
}

func (f *fakeLembreteRepository) Atualizar(
	ctx context.Context,
	lembrete *models.Lembrete,
) error {
	f.atualizarChamado = true
	f.lembreteRecebido = lembrete

	return f.erroAtualizar
}

func (f *fakeLembreteRepository) Excluir(
	ctx context.Context,
	id int64,
) error {
	f.excluirChamado = true
	f.idRecebido = id

	return f.erroExcluir
}

func criarLembreteValido() *models.Lembrete {
	return &models.Lembrete{
		ID:         1,
		Descricao:  "Entregar trabalho de PLP",
		Tipo:       models.TipoEntrega,
		Data:       "2026-08-25",
		Horario:    "14:30",
		Recorrente: false,
	}
}

func TestLembreteServiceSalvar(t *testing.T) {
	repository := &fakeLembreteRepository{}
	service := NewLembreteService(repository)

	lembrete := criarLembreteValido()

	err := service.Salvar(
		context.Background(),
		lembrete,
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

	if repository.lembreteRecebido != lembrete {
		t.Fatal("o lembrete recebido pelo repository é diferente do enviado")
	}
}

func TestLembreteServiceSalvarLembreteNulo(t *testing.T) {
	repository := &fakeLembreteRepository{}
	service := NewLembreteService(repository)

	err := service.Salvar(
		context.Background(),
		nil,
	)

	if !errors.Is(err, ErrLembreteInvalido) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrLembreteInvalido,
		)
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestLembreteServiceSalvarDescricaoVazia(t *testing.T) {
	repository := &fakeLembreteRepository{}
	service := NewLembreteService(repository)

	lembrete := criarLembreteValido()
	lembrete.Descricao = "   "

	err := service.Salvar(
		context.Background(),
		lembrete,
	)

	if !errors.Is(err, ErrDescricaoLembreteObrigatoria) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrDescricaoLembreteObrigatoria,
		)
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestLembreteServiceSalvarDataInvalida(t *testing.T) {
	repository := &fakeLembreteRepository{}
	service := NewLembreteService(repository)

	lembrete := criarLembreteValido()
	lembrete.Data = "25/08/2026"

	err := service.Salvar(
		context.Background(),
		lembrete,
	)

	if !errors.Is(err, ErrDataLembreteInvalida) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrDataLembreteInvalida,
		)
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestLembreteServiceSalvarHorarioInvalido(t *testing.T) {
	repository := &fakeLembreteRepository{}
	service := NewLembreteService(repository)

	lembrete := criarLembreteValido()
	lembrete.Horario = "30:90"

	err := service.Salvar(
		context.Background(),
		lembrete,
	)

	if !errors.Is(err, ErrHorarioLembreteInvalido) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrHorarioLembreteInvalido,
		)
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestLembreteServiceSalvarTipoInvalido(t *testing.T) {
	repository := &fakeLembreteRepository{}
	service := NewLembreteService(repository)

	lembrete := criarLembreteValido()
	lembrete.Tipo = models.TipoLembrete("invalido")

	err := service.Salvar(
		context.Background(),
		lembrete,
	)

	if err == nil {
		t.Fatal("esperava erro para tipo inválido")
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestLembreteServiceBuscarPorID(t *testing.T) {
	lembreteEsperado := criarLembreteValido()

	repository := &fakeLembreteRepository{
		lembreteRetornado: lembreteEsperado,
	}

	service := NewLembreteService(repository)

	lembrete, err := service.BuscarPorID(
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

	if lembrete != lembreteEsperado {
		t.Fatal("lembrete retornado é diferente do esperado")
	}
}

func TestLembreteServiceBuscarPorIDIDInvalido(t *testing.T) {
	repository := &fakeLembreteRepository{}
	service := NewLembreteService(repository)

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

func TestLembreteServiceBuscarTodos(t *testing.T) {
	lembrete := criarLembreteValido()

	repository := &fakeLembreteRepository{
		lembretesRetornados: []models.Lembrete{
			*lembrete,
		},
	}

	service := NewLembreteService(repository)

	lembretes, err := service.BuscarTodos(
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

	if len(lembretes) != 1 {
		t.Fatalf(
			"quantidade de lembretes = %d; esperada 1",
			len(lembretes),
		)
	}
}

func TestLembreteServiceRecorrenciaSemanal(t *testing.T) {
	repository := &fakeLembreteRepository{
		lembretesRetornados: []models.Lembrete{
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

	service := NewLembreteService(repository)

	lembretes, err := service.BuscarTodos(
		context.Background(),
		"2026-08-17",
		"2026-08-31",
	)

	if err != nil {
		t.Fatalf(
			"BuscarTodos() retornou erro inesperado: %v",
			err,
		)
	}

	if len(lembretes) != 3 {
		t.Fatalf(
			"quantidade de ocorrências = %d; esperada 3",
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
				"data da ocorrência %d = %s; esperada %s",
				i,
				lembretes[i].Data,
				data,
			)
		}
	}
}

func TestLembreteServiceBuscarTodosPeriodoInvalido(t *testing.T) {
	repository := &fakeLembreteRepository{}
	service := NewLembreteService(repository)

	_, err := service.BuscarTodos(
		context.Background(),
		"2026-08-31",
		"2026-08-01",
	)

	if !errors.Is(err, ErrPeriodoLembreteInvalido) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrPeriodoLembreteInvalido,
		)
	}

	if repository.buscarTodosChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestLembreteServiceAtualizar(t *testing.T) {
	repository := &fakeLembreteRepository{}
	service := NewLembreteService(repository)

	lembrete := criarLembreteValido()

	err := service.Atualizar(
		context.Background(),
		lembrete,
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

	if repository.lembreteRecebido != lembrete {
		t.Fatal("lembrete recebido pelo repository é diferente")
	}
}

func TestLembreteServiceAtualizarLembreteNulo(t *testing.T) {
	repository := &fakeLembreteRepository{}
	service := NewLembreteService(repository)

	err := service.Atualizar(
		context.Background(),
		nil,
	)

	if !errors.Is(err, ErrLembreteInvalido) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrLembreteInvalido,
		)
	}

	if repository.atualizarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestLembreteServiceAtualizarIDInvalido(t *testing.T) {
	repository := &fakeLembreteRepository{}
	service := NewLembreteService(repository)

	lembrete := criarLembreteValido()
	lembrete.ID = 0

	err := service.Atualizar(
		context.Background(),
		lembrete,
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

func TestLembreteServiceExcluir(t *testing.T) {
	repository := &fakeLembreteRepository{}
	service := NewLembreteService(repository)

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

func TestLembreteServiceExcluirIDInvalido(t *testing.T) {
	repository := &fakeLembreteRepository{}
	service := NewLembreteService(repository)

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