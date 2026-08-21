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
	Lembrete  *repositories.LembreteRepository
}

type Services struct {
	Exemplo   *services.ExemploService
	Meta      *services.MetaService
	Lembrete  *services.LembreteService
}

type Handlers struct {
	Exemplo   *handlers.ExemploHandler
	Meta      *handlers.MetaHandler
	Lembrete  *handlers.LembreteHandler
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

		Lembrete: initializeLembreteRepository(db),
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

		Lembrete: initializeLembreteService(
			appRepositories.Lembrete,
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

		Lembrete: initializeLembreteHandler(
			appServices.Lembrete,
		),
	}
}