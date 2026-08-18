package models

import (
	"testing"
	"time"
)

func TestMetaValidate(t *testing.T) {
	dataInicio := time.Date(
		2026,
		time.August,
		18,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	dataFim := time.Date(
		2026,
		time.August,
		25,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	testes := []struct {
		nome string
		meta Meta
		erro bool
	}{
		{
			nome: "meta válida",
			meta: Meta{
				Status:     StatusCumprida,
				Periodo:    PeriodoSemanal,
				DataInicio: dataInicio,
				DataFim:    dataFim,
			},
			erro: false,
		},
		{
			nome: "status inválido",
			meta: Meta{
				Status:     Status("invalido"),
				Periodo:    PeriodoSemanal,
				DataInicio: dataInicio,
				DataFim:    dataFim,
			},
			erro: true,
		},
		{
			nome: "período inválido",
			meta: Meta{
				Status:     StatusCumprida,
				Periodo:    Periodo("invalido"),
				DataInicio: dataInicio,
				DataFim:    dataFim,
			},
			erro: true,
		},
		{
			nome: "data inicial posterior à final",
			meta: Meta{
				Status:     StatusCumprida,
				Periodo:    PeriodoSemanal,
				DataInicio: dataFim,
				DataFim:    dataInicio,
			},
			erro: true,
		},
		{
			nome: "datas iguais",
			meta: Meta{
				Status:     StatusCumprida,
				Periodo:    PeriodoMensal,
				DataInicio: dataInicio,
				DataFim:    dataInicio,
			},
			erro: false,
		},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			err := tt.meta.Validate()

			if (err != nil) != tt.erro {
				t.Fatalf(
					"Validate() erro = %v; esperado erro = %v",
					err,
					tt.erro,
				)
			}
		})
	}
}

func TestPeriodosValidos(t *testing.T) {
	periodos := []Periodo{
		PeriodoSemanal,
		PeriodoMensal,
		PeriodoAnual,
	}

	for _, periodo := range periodos {
		t.Run(string(periodo), func(t *testing.T) {
			meta := Meta{
				Status:     StatusCumprida,
				Periodo:    periodo,
				DataInicio: time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC),
				DataFim:    time.Date(2026, 8, 25, 0, 0, 0, 0, time.UTC),
			}

			if err := meta.Validate(); err != nil {
				t.Fatalf(
					"período %q deveria ser válido: %v",
					periodo,
					err,
				)
			}
		})
	}
}

func TestStatusValidos(t *testing.T) {
	statuses := []Status{
		StatusCumprida,
		StatusParcialmenteCumprida,
		StatusNaoCumprida,
	}

	for _, status := range statuses {
		t.Run(string(status), func(t *testing.T) {
			meta := Meta{
				Status:     status,
				Periodo:    PeriodoSemanal,
				DataInicio: time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC),
				DataFim:    time.Date(2026, 8, 25, 0, 0, 0, 0, time.UTC),
			}

			if err := meta.Validate(); err != nil {
				t.Fatalf(
					"status %q deveria ser válido: %v",
					status,
					err,
				)
			}
		})
	}
}
