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

type CategoriaHandler struct {
	service *services.CategoriaService
}

func NewCategoriaHandler(
	service *services.CategoriaService,
) *CategoriaHandler {
	return &CategoriaHandler{
		service: service,
	}
}

func (h *CategoriaHandler) BuscarTodos(
	w http.ResponseWriter,
	r *http.Request,
) {
	categorias, err := h.service.BuscarTodos(r.Context())

	if err != nil {
		log.Printf(
			"erro ao buscar categorias: %v",
			err,
		)

		http.Error(
			w,
			"erro ao buscar categorias",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(categorias); err != nil {
		log.Printf(
			"erro ao serializar categorias: %v",
			err,
		)
	}
}

func (h *CategoriaHandler) BuscarPorID(
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

	categoria, err := h.service.BuscarPorID(
		r.Context(),
		id,
	)

	if err != nil {
		if errors.Is(err, repositories.ErrCategoriaNaoEncontrada) {
			http.Error(
				w,
				"categoria não encontrada",
				http.StatusNotFound,
			)

			return
		}

		log.Printf(
			"erro ao buscar categoria: %v",
			err,
		)

		http.Error(
			w,
			"erro ao buscar categoria",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(categoria); err != nil {
		log.Printf(
			"erro ao serializar categoria: %v",
			err,
		)
	}
}

func (h *CategoriaHandler) Criar(
	w http.ResponseWriter,
	r *http.Request,
) {
	var categoria models.Categoria

	err := json.NewDecoder(
		r.Body,
	).Decode(&categoria)

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
		&categoria,
	)

	if err != nil {
		log.Printf(
			"erro ao criar categoria: %v",
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

	if err := json.NewEncoder(w).Encode(categoria); err != nil {
		log.Printf(
			"erro ao serializar categoria criada: %v",
			err,
		)
	}
}

func (h *CategoriaHandler) Atualizar(
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

	var categoria models.Categoria

	err = json.NewDecoder(
		r.Body,
	).Decode(&categoria)

	if err != nil {
		http.Error(
			w,
			"dados inválidos",
			http.StatusBadRequest,
		)

		return
	}

	categoria.ID = id

	err = h.service.Atualizar(
		r.Context(),
		&categoria,
	)

	if err != nil {
		if errors.Is(err, repositories.ErrCategoriaNaoEncontrada) {
			http.Error(
				w,
				"categoria não encontrada",
				http.StatusNotFound,
			)

			return
		}

		log.Printf(
			"erro ao atualizar categoria: %v",
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

	if err := json.NewEncoder(w).Encode(categoria); err != nil {
		log.Printf(
			"erro ao serializar categoria atualizada: %v",
			err,
		)
	}
}

func (h *CategoriaHandler) Excluir(
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
		if errors.Is(err, repositories.ErrCategoriaNaoEncontrada) {
			http.Error(
				w,
				"categoria não encontrada",
				http.StatusNotFound,
			)

			return
		}

		log.Printf(
			"erro ao excluir categoria: %v",
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
