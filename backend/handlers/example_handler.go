package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"plp-planner/services"
)

type ExemploHandler struct {
	service *services.ExemploService
}

func NewExemploHandler(
	service *services.ExemploService,
) *ExemploHandler {
	return &ExemploHandler{
		service: service,
	}
}

func (h *ExemploHandler) BuscarTodos(
	w http.ResponseWriter,
	r *http.Request,
) {
	exemplos, err := h.service.BuscarTodos(r.Context())
	if err != nil {
		log.Printf("Erro ao buscar exemplos: %v", err)

		http.Error(
			w,
			"Erro ao buscar exemplos",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exemplos)
}