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

type TarefaHandler struct {
	service *services.TarefaService
}

func NewTarefaHandler(
	service *services.TarefaService,
) *TarefaHandler {
	return &TarefaHandler{
		service: service,
	}
}

func (h *TarefaHandler) BuscarTodos(
	w http.ResponseWriter,
	r *http.Request,
) {
	query := r.URL.Query()

	data := query.Get("data")
	categoriaID := query.Get("categoria_id")

	tarefas, err := h.service.BuscarTodos(
		r.Context(),
		data,
		categoriaID,
	)

	if err != nil {
		log.Printf(
			"erro ao buscar tarefas: %v",
			err,
		)

		http.Error(
			w,
			"erro ao buscar tarefas",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tarefas); err != nil {
		log.Printf(
			"erro ao serializar tarefas: %v",
			err,
		)
	}
}

func (h *TarefaHandler) BuscarPorID(
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

	tarefa, err := h.service.BuscarPorID(
		r.Context(),
		id,
	)

	if err != nil {
		if errors.Is(err, repositories.ErrTarefaNaoEncontrada) {
			http.Error(
				w,
				"tarefa não encontrada",
				http.StatusNotFound,
			)

			return
		}

		log.Printf(
			"erro ao buscar tarefa: %v",
			err,
		)

		http.Error(
			w,
			"erro ao buscar tarefa",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tarefa); err != nil {
		log.Printf(
			"erro ao serializar tarefa: %v",
			err,
		)
	}
}

func (h *TarefaHandler) Criar(
	w http.ResponseWriter,
	r *http.Request,
) {
	var tarefa models.Tarefa

	err := json.NewDecoder(
		r.Body,
	).Decode(&tarefa)

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
		&tarefa,
	)

	if err != nil {
		log.Printf(
			"erro ao criar tarefa: %v",
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

	if err := json.NewEncoder(w).Encode(tarefa); err != nil {
		log.Printf(
			"erro ao serializar tarefa criada: %v",
			err,
		)
	}
}

func (h *TarefaHandler) Atualizar(
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

	var tarefa models.Tarefa

	err = json.NewDecoder(
		r.Body,
	).Decode(&tarefa)

	if err != nil {
		http.Error(
			w,
			"dados inválidos",
			http.StatusBadRequest,
		)

		return
	}

	tarefa.ID = id

	err = h.service.Atualizar(
		r.Context(),
		&tarefa,
	)

	if err != nil {
		if errors.Is(err, repositories.ErrTarefaNaoEncontrada) {
			http.Error(
				w,
				"tarefa não encontrada",
				http.StatusNotFound,
			)

			return
		}

		log.Printf(
			"erro ao atualizar tarefa: %v",
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

	if err := json.NewEncoder(w).Encode(tarefa); err != nil {
		log.Printf(
			"erro ao serializar tarefa atualizada: %v",
			err,
		)
	}
}

func (h *TarefaHandler) AtualizarStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := obterIDStatus(r.URL.Path)

	if err != nil {
		http.Error(
			w,
			"id inválido",
			http.StatusBadRequest,
		)

		return
	}

	var request struct {
		Status models.StatusTarefa `json:"status"`
	}

	err = json.NewDecoder(
		r.Body,
	).Decode(&request)

	if err != nil {
		http.Error(
			w,
			"dados inválidos",
			http.StatusBadRequest,
		)

		return
	}

	err = h.service.AtualizarStatus(
		r.Context(),
		id,
		request.Status,
	)

	if err != nil {
		if errors.Is(err, repositories.ErrTarefaNaoEncontrada) {
			http.Error(
				w,
				"tarefa não encontrada",
				http.StatusNotFound,
			)

			return
		}

		log.Printf(
			"erro ao atualizar status: %v",
			err,
		)

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TarefaHandler) Excluir(
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
		if errors.Is(err, repositories.ErrTarefaNaoEncontrada) {
			http.Error(
				w,
				"tarefa não encontrada",
				http.StatusNotFound,
			)

			return
		}

		log.Printf(
			"erro ao excluir tarefa: %v",
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
