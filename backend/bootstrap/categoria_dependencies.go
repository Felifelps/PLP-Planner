package bootstrap

import (
	"plp-planner/handlers"
	"plp-planner/repositories"
	"plp-planner/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

func initializeCategoriaRepository(
	db *pgxpool.Pool,
) *repositories.CategoriaRepository {
	return repositories.NewCategoriaRepository(db)
}

func initializeCategoriaService(
	repository *repositories.CategoriaRepository,
) *services.CategoriaService {
	return services.NewCategoriaService(repository)
}

func initializeCategoriaHandler(
	service *services.CategoriaService,
) *handlers.CategoriaHandler {
	return handlers.NewCategoriaHandler(service)
}
