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
	initializeCategoriaRoutes(router, dependencies.Handlers)
	initializeTarefaRoutes(router, dependencies.Handlers)

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

func initializeCategoriaRoutes(
	router *http.ServeMux,
	appHandlers *Handlers,
) {
	router.HandleFunc(
		"GET /categorias",
		appHandlers.Categoria.BuscarTodos,
	)

	router.HandleFunc(
		"GET /categorias/{id}",
		appHandlers.Categoria.BuscarPorID,
	)

	router.HandleFunc(
		"POST /categorias",
		appHandlers.Categoria.Criar,
	)

	router.HandleFunc(
		"PUT /categorias/{id}",
		appHandlers.Categoria.Atualizar,
	)

	router.HandleFunc(
		"DELETE /categorias/{id}",
		appHandlers.Categoria.Excluir,
	)
}

func initializeTarefaRoutes(
	router *http.ServeMux,
	appHandlers *Handlers,
) {
	router.HandleFunc(
		"GET /tarefas",
		appHandlers.Tarefa.BuscarTodos,
	)

	router.HandleFunc(
		"GET /tarefas/{id}",
		appHandlers.Tarefa.BuscarPorID,
	)

	router.HandleFunc(
		"POST /tarefas",
		appHandlers.Tarefa.Criar,
	)

	router.HandleFunc(
		"PUT /tarefas/{id}",
		appHandlers.Tarefa.Atualizar,
	)

	router.HandleFunc(
		"PATCH /tarefas/{id}/status",
		appHandlers.Tarefa.AtualizarStatus,
	)

	router.HandleFunc(
		"DELETE /tarefas/{id}",
		appHandlers.Tarefa.Excluir,
	)
}
