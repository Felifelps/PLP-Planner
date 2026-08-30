package bootstrap

import (
	"plp-planner/handlers"
	"plp-planner/repositories"
	"plp-planner/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	Exemplo   *repositories.ExemploRepository
	Meta      *repositories.MetaRepository
	Categoria *repositories.CategoriaRepository
	Tarefa    *repositories.TarefaRepository
	Lembrete  *repositories.LembreteRepository
}

type Services struct {
	Exemplo   *services.ExemploService
	Meta      *services.MetaService
	Categoria *services.CategoriaService
	Tarefa    *services.TarefaService
	Lembrete  *services.LembreteService
	Relatorio *services.RelatorioService
}

type Handlers struct {
	Exemplo   *handlers.ExemploHandler
	Meta      *handlers.MetaHandler
	Categoria *handlers.CategoriaHandler
	Tarefa    *handlers.TarefaHandler
	Lembrete  *handlers.LembreteHandler
	Relatorio *handlers.RelatorioHandler
}

type Dependencies struct {
	Repositories *Repositories
	Services     *Services
	Handlers     *Handlers
}

func InitializeDependencies(
	db *pgxpool.Pool,
) *Dependencies {
	appRepositories := initializeRepositories(db)
	appServices := initializeServices(appRepositories)
	appHandlers := initializeHandlers(appServices)

	return &Dependencies{
		Repositories: appRepositories,
		Services:     appServices,
		Handlers:     appHandlers,
	}
}

func initializeRepositories(
	db *pgxpool.Pool,
) *Repositories {
	return &Repositories{
		Exemplo: repositories.NewExemploRepository(db),

		Meta: initializeMetaRepository(db),

		Categoria: initializeCategoriaRepository(db),
		Tarefa:    initializeTarefaRepository(db),
		Lembrete:  initializeLembreteRepository(db),
	}
}

func initializeServices(
	appRepositories *Repositories,
) *Services {
	return &Services{
		Exemplo: services.NewExemploService(
			appRepositories.Exemplo,
		),

		Meta: initializeMetaService(appRepositories.Meta),

		Categoria: initializeCategoriaService(appRepositories.Categoria),
		Tarefa:    initializeTarefaService(appRepositories.Tarefa),
		Lembrete: initializeLembreteService(
			appRepositories.Lembrete,
		),
		Relatorio: initializeRelatorioService(
			appRepositories.Meta,
			appRepositories.Tarefa,
		),
	}
}

func initializeHandlers(
	appServices *Services,
) *Handlers {
	return &Handlers{
		Exemplo: handlers.NewExemploHandler(
			appServices.Exemplo,
		),

		Meta: initializeMetaHandler(appServices.Meta),

		Categoria: initializeCategoriaHandler(appServices.Categoria),
		Tarefa:    initializeTarefaHandler(appServices.Tarefa),
		Lembrete: initializeLembreteHandler(
			appServices.Lembrete,
		),
		Relatorio: initializeRelatorioHandler(
			appServices.Relatorio,
		),
	}
}