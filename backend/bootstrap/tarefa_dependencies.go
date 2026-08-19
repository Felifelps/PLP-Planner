package bootstrap

import (
	"plp-planner/handlers"
	"plp-planner/repositories"
	"plp-planner/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

func initializeTarefaRepository(
	db *pgxpool.Pool,
) *repositories.TarefaRepository {
	return repositories.NewTarefaRepository(db)
}

func initializeTarefaService(
	repository *repositories.TarefaRepository,
) *services.TarefaService {
	return services.NewTarefaService(repository)
}

func initializeTarefaHandler(
	service *services.TarefaService,
) *handlers.TarefaHandler {
	return handlers.NewTarefaHandler(service)
}
