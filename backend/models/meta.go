package models

import (
	"errors"
	"time"
)

type Status string

const (
	StatusCumprida             Status = "cumprida"
	StatusParcialmenteCumprida Status = "parcialmente cumprida"
	StatusNaoCumprida          Status = "não cumprida"
)

type Periodo string

const (
	PeriodoDiario Periodo = "diario"
	PeriodoSemanal Periodo = "semanal"
	PeriodoMensal  Periodo = "mensal"
	PeriodoAnual   Periodo = "anual"
)

type Meta struct {
	ID          int64     `json:"id"`
	Nome        string    `json:"nome"`
	Descricao   string    `json:"descricao"`
	CategoriaID int64     `json:"categoria_id"`
	Status      Status    `json:"status"`
	Periodo     Periodo   `json:"periodo"`
	DataInicio  time.Time `json:"data_inicio"`
	DataFim     time.Time `json:"data_fim"`
}

func (m *Meta) Validate() error {
	if !StatusValido(m.Status) {
		return errors.New("status inválido")
	}

	if !PeriodoValido(m.Periodo) {
		return errors.New("período inválido")
	}

	if m.DataInicio.After(m.DataFim) {
		return errors.New(
			"a data de início não pode ser posterior à data de fim",
		)
	}

	return nil
}

func StatusValido(status Status) bool {
	switch status {
	case StatusCumprida,
		StatusParcialmenteCumprida,
		StatusNaoCumprida:
		return true
	default:
		return false
	}
}

func PeriodoValido(periodo Periodo) bool {
	switch periodo {
	case PeriodoDiario, 
		PeriodoSemanal,
		PeriodoMensal,
		PeriodoAnual:
		return true
	default:
		return false
	}
}
