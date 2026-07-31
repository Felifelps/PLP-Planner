package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"plp-planner/models"
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

	metas, err := h.service.BuscarTodos(
		r.Context(),
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

	json.NewEncoder(
		w,
	).Encode(metas)
}

func (h *MetaHandler) BuscarPorID(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := obterID(
		r.URL.Path,
	)

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

	json.NewEncoder(
		w,
	).Encode(meta)
}

func (h *MetaHandler) Criar(
	w http.ResponseWriter,
	r *http.Request,
) {

	var meta models.Meta

	err := json.NewDecoder(
		r.Body,
	).Decode(
		&meta,
	)

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

	w.WriteHeader(
		http.StatusCreated,
	)

	json.NewEncoder(
		w,
	).Encode(meta)
}

func (h *MetaHandler) Atualizar(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := obterID(
		r.URL.Path,
	)

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
	).Decode(
		&meta,
	)

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

	json.NewEncoder(
		w,
	).Encode(meta)
}

func (h *MetaHandler) AtualizarStatus(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := obterID(
		r.URL.Path,
	)

	if err != nil {
		http.Error(
			w,
			"id inválido",
			http.StatusBadRequest,
		)

		return
	}

	var request struct {
		Status string `json:"status"`
	}

	err = json.NewDecoder(
		r.Body,
	).Decode(
		&request,
	)

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

	w.WriteHeader(
		http.StatusNoContent,
	)
}

func (h *MetaHandler) Excluir(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := obterID(
		r.URL.Path,
	)

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
		log.Printf(
			"erro ao excluir meta: %v",
			err,
		)

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	w.WriteHeader(
		http.StatusNoContent,
	)
}

func obterID(
	path string,
) (int64, error) {

	partes := strings.Split(
		strings.Trim(
			path,
			"/",
		),
		"/",
	)

	id := partes[len(partes)-1]

	return strconv.ParseInt(
		id,
		10,
		64,
	)
}
