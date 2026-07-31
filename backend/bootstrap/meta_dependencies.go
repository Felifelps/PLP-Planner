package bootstrap

import (
	"plp-planner/handlers"
	"plp-planner/repositories"
	"plp-planner/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

func initializeMetaRepository(
	db *pgxpool.Pool,
) *repositories.MetaRepository {

	return repositories.NewMetaRepository(
		db,
	)
}

func initializeMetaService(
	repository *repositories.MetaRepository,
) *services.MetaService {

	return services.NewMetaService(
		repository,
	)
}

func initializeMetaHandler(
	service *services.MetaService,
) *handlers.MetaHandler {

	return handlers.NewMetaHandler(
		service,
	)
}
