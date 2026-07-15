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