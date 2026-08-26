package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"plp-planner/models"
)

type fakeTarefaRepository struct {
	salvarChamado          bool
	buscarTodosChamado     bool
	buscarPorIDChamado     bool
	atualizarChamado       bool
	atualizarStatusChamado bool
	excluirChamado         bool

	tarefaRecebida *models.Tarefa
	idRecebido     int64
	statusRecebido models.StatusTarefa

	tarefasRetornadas []models.Tarefa
	tarefaRetornada   *models.Tarefa

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
	f.salvarChamado = true
	f.tarefaRecebida = tarefa
	return f.erroSalvar
}

func (f *fakeTarefaRepository) BuscarTodos(
	ctx context.Context,
	data string,
	categoriaID string,
) ([]models.Tarefa, error) {
	f.buscarTodosChamado = true
	return f.tarefasRetornadas, f.erroBuscarTodos
}

func (f *fakeTarefaRepository) BuscarPorID(
	ctx context.Context,
	id int64,
) (*models.Tarefa, error) {
	f.buscarPorIDChamado = true
	f.idRecebido = id
	return f.tarefaRetornada, f.erroBuscarPorID
}

func (f *fakeTarefaRepository) Atualizar(
	ctx context.Context,
	tarefa *models.Tarefa,
) error {
	f.atualizarChamado = true
	f.tarefaRecebida = tarefa
	return f.erroAtualizar
}

func (f *fakeTarefaRepository) AtualizarStatus(
	ctx context.Context,
	id int64,
	status models.StatusTarefa,
) error {
	f.atualizarStatusChamado = true
	f.idRecebido = id
	f.statusRecebido = status
	return f.erroAtualizarStatus
}

func (f *fakeTarefaRepository) Excluir(
	ctx context.Context,
	id int64,
) error {
	f.excluirChamado = true
	f.idRecebido = id
	return f.erroExcluir
}

func criarTarefaValida() *models.Tarefa {
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

func TestTarefaServiceSalvar(t *testing.T) {
	repository := &fakeTarefaRepository{}
	service := NewTarefaService(repository)

	tarefa := criarTarefaValida()

	err := service.Salvar(
		context.Background(),
		tarefa,
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

	if repository.tarefaRecebida != tarefa {
		t.Fatal("a tarefa recebida pelo repository é diferente da enviada")
	}
}

func TestTarefaServiceSalvarTarefaNula(t *testing.T) {
	repository := &fakeTarefaRepository{}
	service := NewTarefaService(repository)

	err := service.Salvar(
		context.Background(),
		nil,
	)

	if !errors.Is(err, ErrTarefaInvalida) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrTarefaInvalida,
		)
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestTarefaServiceSalvarDescricaoVazia(t *testing.T) {
	repository := &fakeTarefaRepository{}
	service := NewTarefaService(repository)

	tarefa := criarTarefaValida()
	tarefa.Descricao = "   "

	err := service.Salvar(
		context.Background(),
		tarefa,
	)

	if !errors.Is(err, ErrDescricaoObrigatoria) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrDescricaoObrigatoria,
		)
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestTarefaServiceSalvarCategoriaInvalida(t *testing.T) {
	repository := &fakeTarefaRepository{}
	service := NewTarefaService(repository)

	tarefa := criarTarefaValida()
	tarefa.CategoriaID = 0

	err := service.Salvar(
		context.Background(),
		tarefa,
	)

	if !errors.Is(err, ErrTarefaCategoriaInvalida) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrTarefaCategoriaInvalida,
		)
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestTarefaServiceSalvarHorarioETurnoJuntos(t *testing.T) {
	repository := &fakeTarefaRepository{}
	service := NewTarefaService(repository)

	horario := "10:00"
	duracao := models.Duracao1Hora

	tarefa := criarTarefaValida()
	tarefa.HorarioInicio = &horario
	tarefa.Duracao = &duracao

	err := service.Salvar(
		context.Background(),
		tarefa,
	)

	if err == nil {
		t.Fatal("esperava erro ao informar horário e turno juntos")
	}

	if repository.salvarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestTarefaServiceBuscarPorID(t *testing.T) {
	tarefaEsperada := criarTarefaValida()

	repository := &fakeTarefaRepository{
		tarefaRetornada: tarefaEsperada,
	}

	service := NewTarefaService(repository)

	tarefa, err := service.BuscarPorID(
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

	if tarefa != tarefaEsperada {
		t.Fatal("tarefa retornada é diferente da esperada")
	}
}

func TestTarefaServiceBuscarPorIDIDInvalido(t *testing.T) {
	repository := &fakeTarefaRepository{}
	service := NewTarefaService(repository)

	_, err := service.BuscarPorID(
		context.Background(),
		0,
	)

	if !errors.Is(err, ErrTarefaIDInvalido) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrTarefaIDInvalido,
		)
	}

	if repository.buscarPorIDChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestTarefaServiceBuscarTodos(t *testing.T) {
	tarefasEsperadas := []models.Tarefa{
		*criarTarefaValida(),
	}

	repository := &fakeTarefaRepository{
		tarefasRetornadas: tarefasEsperadas,
	}

	service := NewTarefaService(repository)

	tarefas, err := service.BuscarTodos(
		context.Background(),
		"2026-08-19",
		"1",
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

	if len(tarefas) != 1 {
		t.Fatalf(
			"quantidade de tarefas = %d; esperada 1",
			len(tarefas),
		)
	}
}

func TestTarefaServiceAtualizar(t *testing.T) {
	repository := &fakeTarefaRepository{}
	service := NewTarefaService(repository)

	tarefa := criarTarefaValida()

	err := service.Atualizar(
		context.Background(),
		tarefa,
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

func TestTarefaServiceAtualizarIDInvalido(t *testing.T) {
	repository := &fakeTarefaRepository{}
	service := NewTarefaService(repository)

	tarefa := criarTarefaValida()
	tarefa.ID = 0

	err := service.Atualizar(
		context.Background(),
		tarefa,
	)

	if !errors.Is(err, ErrTarefaIDInvalido) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrTarefaIDInvalido,
		)
	}

	if repository.atualizarChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestTarefaServiceAtualizarStatus(t *testing.T) {
	repository := &fakeTarefaRepository{}
	service := NewTarefaService(repository)

	err := service.AtualizarStatus(
		context.Background(),
		1,
		models.StatusAdiada,
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

	if repository.statusRecebido != models.StatusAdiada {
		t.Fatalf(
			"status recebido = %q; esperado %q",
			repository.statusRecebido,
			models.StatusAdiada,
		)
	}
}

func TestTarefaServiceAtualizarStatusIDInvalido(t *testing.T) {
	repository := &fakeTarefaRepository{}
	service := NewTarefaService(repository)

	err := service.AtualizarStatus(
		context.Background(),
		0,
		models.StatusExecutada,
	)

	if !errors.Is(err, ErrTarefaIDInvalido) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrTarefaIDInvalido,
		)
	}

	if repository.atualizarStatusChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestTarefaServiceAtualizarStatusInvalido(t *testing.T) {
	repository := &fakeTarefaRepository{}
	service := NewTarefaService(repository)

	err := service.AtualizarStatus(
		context.Background(),
		1,
		models.StatusTarefa("invalido"),
	)

	if !errors.Is(err, ErrTarefaStatusInvalido) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrTarefaStatusInvalido,
		)
	}

	if repository.atualizarStatusChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestTarefaServiceExcluir(t *testing.T) {
	repository := &fakeTarefaRepository{}
	service := NewTarefaService(repository)

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

func TestTarefaServiceExcluirIDInvalido(t *testing.T) {
	repository := &fakeTarefaRepository{}
	service := NewTarefaService(repository)

	err := service.Excluir(
		context.Background(),
		0,
	)

	if !errors.Is(err, ErrTarefaIDInvalido) {
		t.Fatalf(
			"erro = %v; esperado %v",
			err,
			ErrTarefaIDInvalido,
		)
	}

	if repository.excluirChamado {
		t.Fatal("repository não deveria ser chamado")
	}
}
