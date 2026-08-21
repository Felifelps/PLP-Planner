package bootstrap

import (
	"plp-planner/handlers"
	"plp-planner/repositories"
	"plp-planner/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

func initializeLembreteRepository(
	db *pgxpool.Pool,
) *repositories.LembreteRepository {
	return repositories.NewLembreteRepository(db)
}

func initializeLembreteService(
	repository *repositories.LembreteRepository,
) *services.LembreteService {
	return services.NewLembreteService(repository)
}

func initializeLembreteHandler(
	service *services.LembreteService,
) *handlers.LembreteHandler {
	return handlers.NewLembreteHandler(service)
}