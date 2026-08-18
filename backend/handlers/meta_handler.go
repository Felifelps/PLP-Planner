package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"plp-planner/models"
	"plp-planner/repositories"
	"plp-planner/services"
)

type MetaHandler struct {
	service *services.MetaService
}

func NewMetaHandler(
	service *services.MetaService,
) *MetaHandler {
	return &MetaHandler{
		service: service,
	}
}

func (h *MetaHandler) BuscarTodos(
	w http.ResponseWriter,
	r *http.Request,
) {
	query := r.URL.Query()

	dataInicio := query.Get("data_inicio")
	dataFim := query.Get("data_fim")

	metas, err := h.service.BuscarTodos(
		r.Context(),
		dataInicio,
		dataFim,
	)

	if err != nil {
		log.Printf(
			"erro ao buscar metas: %v",
			err,
		)

		http.Error(
			w,
			"erro ao buscar metas",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(metas); err != nil {
		log.Printf(
			"erro ao serializar metas: %v",
			err,
		)
	}
}

func (h *MetaHandler) BuscarPorID(
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

	meta, err := h.service.BuscarPorID(
		r.Context(),
		id,
	)

	if err != nil {
		if errors.Is(err, repositories.ErrMetaNaoEncontrada) {
			http.Error(
				w,
				"meta não encontrada",
				http.StatusNotFound,
			)

			return
		}

		log.Printf(
			"erro ao buscar meta: %v",
			err,
		)

		http.Error(
			w,
			"erro ao buscar meta",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(meta); err != nil {
		log.Printf(
			"erro ao serializar meta: %v",
			err,
		)
	}
}

func (h *MetaHandler) Criar(
	w http.ResponseWriter,
	r *http.Request,
) {
	var meta models.Meta

	err := json.NewDecoder(
		r.Body,
	).Decode(&meta)

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
		&meta,
	)

	if err != nil {
		log.Printf(
			"erro ao criar meta: %v",
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

	if err := json.NewEncoder(w).Encode(meta); err != nil {
		log.Printf(
			"erro ao serializar meta criada: %v",
			err,
		)
	}
}

func (h *MetaHandler) Atualizar(
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

	var meta models.Meta

	err = json.NewDecoder(
		r.Body,
	).Decode(&meta)

	if err != nil {
		http.Error(
			w,
			"dados inválidos",
			http.StatusBadRequest,
		)

		return
	}

	meta.ID = id

	err = h.service.Atualizar(
		r.Context(),
		&meta,
	)

	if err != nil {
		if errors.Is(err, repositories.ErrMetaNaoEncontrada) {
			http.Error(
				w,
				"meta não encontrada",
				http.StatusNotFound,
			)

			return
		}

		log.Printf(
			"erro ao atualizar meta: %v",
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

	if err := json.NewEncoder(w).Encode(meta); err != nil {
		log.Printf(
			"erro ao serializar meta atualizada: %v",
			err,
		)
	}
}

func (h *MetaHandler) AtualizarStatus(
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
		Status models.Status `json:"status"`
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
		if errors.Is(err, repositories.ErrMetaNaoEncontrada) {
			http.Error(
				w,
				"meta não encontrada",
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

func (h *MetaHandler) Excluir(
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
		if errors.Is(err, repositories.ErrMetaNaoEncontrada) {
			http.Error(
				w,
				"meta não encontrada",
				http.StatusNotFound,
			)

			return
		}

		log.Printf(
			"erro ao excluir meta: %v",
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

func obterID(path string) (int64, error) {
	partes := strings.Split(
		strings.Trim(path, "/"),
		"/",
	)

	if len(partes) == 0 || partes[len(partes)-1] == "" {
		return 0, errors.New("id inválido")
	}

	id := partes[len(partes)-1]

	return strconv.ParseInt(
		id,
		10,
		64,
	)
}

func obterIDStatus(path string) (int64, error) {
	partes := strings.Split(
		strings.Trim(path, "/"),
		"/",
	)

	if len(partes) < 3 {
		return 0, errors.New("id inválido")
	}

	if partes[len(partes)-1] != "status" {
		return 0, errors.New("id inválido")
	}

	id := partes[len(partes)-2]

	return strconv.ParseInt(
		id,
		10,
		64,
	)
}
