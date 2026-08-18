package bootstrap

import (
	"net/http"

	"plp-planner/handlers"
)

func InitializeRouter(
	dependencies *Dependencies,
) *http.ServeMux {
	router := http.NewServeMux()

	initializeStatusRoutes(router)
	initializeExemploRoutes(router, dependencies.Handlers)
	initializeMetaRoutes(router, dependencies.Handlers)

	return router
}

func initializeStatusRoutes(
	router *http.ServeMux,
) {
	router.HandleFunc(
		"GET /",
		handlers.Status,
	)
}

func initializeExemploRoutes(
	router *http.ServeMux,
	appHandlers *Handlers,
) {
	router.HandleFunc(
		"GET /exemplos",
		appHandlers.Exemplo.BuscarTodos,
	)
}

func initializeMetaRoutes(
	router *http.ServeMux,
	appHandlers *Handlers,
) {
	router.HandleFunc(
		"GET /metas",
		appHandlers.Meta.BuscarTodos,
	)

	router.HandleFunc(
		"GET /metas/{id}",
		appHandlers.Meta.BuscarPorID,
	)

	router.HandleFunc(
		"POST /metas",
		appHandlers.Meta.Criar,
	)

	router.HandleFunc(
		"PUT /metas/{id}",
		appHandlers.Meta.Atualizar,
	)

	router.HandleFunc(
		"PATCH /metas/{id}/status",
		appHandlers.Meta.AtualizarStatus,
	)

	router.HandleFunc(
		"DELETE /metas/{id}",
		appHandlers.Meta.Excluir,
	)
}
