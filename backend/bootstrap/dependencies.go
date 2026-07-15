package bootstrap

import (
	"plp-planner/handlers"
	"plp-planner/repositories"
	"plp-planner/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	Exemplo *repositories.ExemploRepository
}

type Services struct {
	Exemplo *services.ExemploService
}

type Handlers struct {
	Exemplo *handlers.ExemploHandler
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
	}
}

func initializeServices(
	appRepositories *Repositories,
) *Services {
	return &Services{
		Exemplo: services.NewExemploService(
			appRepositories.Exemplo,
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
	}
}