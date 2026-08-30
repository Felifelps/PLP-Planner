package bootstrap

import (
	"plp-planner/handlers"
	"plp-planner/services"
)

func initializeRelatorioService(
	metaRepository services.RelatorioMetaRepository,
	tarefaRepository services.RelatorioTarefaRepository,
) *services.RelatorioService {
	return services.NewRelatorioService(metaRepository, tarefaRepository)
}

func initializeRelatorioHandler(
	service *services.RelatorioService,
) *handlers.RelatorioHandler {
	return handlers.NewRelatorioHandler(service)
}
