package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"plp-planner/services"
)

type RelatorioHandler struct {
	service *services.RelatorioService
}

func NewRelatorioHandler(
	service *services.RelatorioService,
) *RelatorioHandler {
	return &RelatorioHandler{
		service: service,
	}
}

func (h *RelatorioHandler) Gerar(
	w http.ResponseWriter,
	r *http.Request,
) {
	query := r.URL.Query()

	dataInicio := query.Get("data_inicio")
	dataFim := query.Get("data_fim")

	relatorio, err := h.service.Gerar(
		r.Context(),
		dataInicio,
		dataFim,
	)

	if err != nil {
		if errors.Is(err, services.ErrPeriodoRelatorioObrigatorio) ||
			errors.Is(err, services.ErrPeriodoRelatorioInvalido) {

			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)

			return
		}

		log.Printf(
			"erro ao gerar relatório: %v",
			err,
		)

		http.Error(
			w,
			"erro ao gerar relatório",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(relatorio); err != nil {
		log.Printf(
			"erro ao serializar relatório: %v",
			err,
		)
	}
}
