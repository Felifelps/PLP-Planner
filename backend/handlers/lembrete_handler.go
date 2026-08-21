package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"plp-planner/models"
	"plp-planner/repositories"
	"plp-planner/services"
)

type LembreteHandler struct {
	service *services.LembreteService
}

func NewLembreteHandler(
	service *services.LembreteService,
) *LembreteHandler {
	return &LembreteHandler{
		service: service,
	}
}

func (h *LembreteHandler) BuscarTodos(
	w http.ResponseWriter,
	r *http.Request,
) {
	query := r.URL.Query()

	dataInicio := query.Get("data_inicio")
	dataFim := query.Get("data_fim")

	lembretes, err := h.service.BuscarTodos(
		r.Context(),
		dataInicio,
		dataFim,
	)

	if err != nil {
		if errors.Is(err, services.ErrDataLembreteInvalida) ||
			errors.Is(err, services.ErrPeriodoLembreteInvalido) {

			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)

			return
		}

		log.Printf(
			"erro ao buscar lembretes: %v",
			err,
		)

		http.Error(
			w,
			"erro ao buscar lembretes",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(lembretes); err != nil {
		log.Printf(
			"erro ao serializar lembretes: %v",
			err,
		)
	}
}

func (h *LembreteHandler) BuscarPorID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := obterID(r.URL.Path)

	if err != nil {
		http.Error(
			w,
			"id inválido",
			http.StatusBadRequest,
		)

		return
	}

	lembrete, err := h.service.BuscarPorID(
		r.Context(),
		id,
	)

	if err != nil {
		if errors.Is(
			err,
			repositories.ErrLembreteNaoEncontrado,
		) {
			http.Error(
				w,
				"lembrete não encontrado",
				http.StatusNotFound,
			)

			return
		}

		if errors.Is(err, services.ErrIDInvalido) {
			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)

			return
		}

		log.Printf(
			"erro ao buscar lembrete: %v",
			err,
		)

		http.Error(
			w,
			"erro ao buscar lembrete",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(lembrete); err != nil {
		log.Printf(
			"erro ao serializar lembrete: %v",
			err,
		)
	}
}

func (h *LembreteHandler) Criar(
	w http.ResponseWriter,
	r *http.Request,
) {
	var lembrete models.Lembrete

	err := json.NewDecoder(
		r.Body,
	).Decode(&lembrete)

	if err != nil {
		http.Error(
			w,
			"dados inválidos",
			http.StatusBadRequest,
		)

		return
	}

	err = h.service.Salvar(
		r.Context(),
		&lembrete,
	)

	if err != nil {
		log.Printf(
			"erro ao criar lembrete: %v",
			err,
		)

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(lembrete); err != nil {
		log.Printf(
			"erro ao serializar lembrete criado: %v",
			err,
		)
	}
}

func (h *LembreteHandler) Atualizar(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := obterID(r.URL.Path)

	if err != nil {
		http.Error(
			w,
			"id inválido",
			http.StatusBadRequest,
		)

		return
	}

	var lembrete models.Lembrete

	err = json.NewDecoder(
		r.Body,
	).Decode(&lembrete)

	if err != nil {
		http.Error(
			w,
			"dados inválidos",
			http.StatusBadRequest,
		)

		return
	}

	lembrete.ID = id

	err = h.service.Atualizar(
		r.Context(),
		&lembrete,
	)

	if err != nil {
		if errors.Is(
			err,
			repositories.ErrLembreteNaoEncontrado,
		) {
			http.Error(
				w,
				"lembrete não encontrado",
				http.StatusNotFound,
			)

			return
		}

		log.Printf(
			"erro ao atualizar lembrete: %v",
			err,
		)

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(lembrete); err != nil {
		log.Printf(
			"erro ao serializar lembrete atualizado: %v",
			err,
		)
	}
}

func (h *LembreteHandler) Excluir(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := obterID(r.URL.Path)

	if err != nil {
		http.Error(
			w,
			"id inválido",
			http.StatusBadRequest,
		)

		return
	}

	err = h.service.Excluir(
		r.Context(),
		id,
	)

	if err != nil {
		if errors.Is(
			err,
			repositories.ErrLembreteNaoEncontrado,
		) {
			http.Error(
				w,
				"lembrete não encontrado",
				http.StatusNotFound,
			)

			return
		}

		log.Printf(
			"erro ao excluir lembrete: %v",
			err,
		)

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}