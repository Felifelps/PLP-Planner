package models

import "time"

const (
	StatusCumprida             = "cumprida"
	StatusParcialmenteCumprida = "parcialmente cumprida"
	StatusNaoCumprida          = "não cumprida"
)

type Meta struct {
	ID          int64     `json:"id"`
	Nome        string    `json:"nome"`
	Descricao   string    `json:"descricao"`
	CategoriaID int64     `json:"categoria_id"`
	Status      string    `json:"status"`
	DataInicio  time.Time `json:"data_inicio"`
	DataFim     time.Time `json:"data_fim"`
}
