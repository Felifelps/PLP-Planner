package models

import "testing"

func TestLembreteValidate(t *testing.T) {
	testes := []struct {
		nome     string
		lembrete Lembrete
		erro     bool
	}{
		{
			nome: "lembrete válido",
			lembrete: Lembrete{
				Descricao:  "Entregar trabalho de PLP",
				Tipo:       TipoEntrega,
				Data:       "2026-08-25",
				Horario:    "14:30",
				Recorrente: false,
			},
			erro: false,
		},
		{
			nome: "tipo inválido",
			lembrete: Lembrete{
				Descricao:  "Lembrete inválido",
				Tipo:       TipoLembrete("invalido"),
				Data:       "2026-08-25",
				Horario:    "14:30",
				Recorrente: false,
			},
			erro: true,
		},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			err := tt.lembrete.Validate()

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

func TestTiposLembreteValidos(t *testing.T) {
	tipos := []TipoLembrete{
		TipoReuniao,
		TipoLigacao,
		TipoCompra,
		TipoEstudo,
		TipoExercicio,
		TipoEntrega,
	}

	for _, tipo := range tipos {
		t.Run(string(tipo), func(t *testing.T) {
			lembrete := Lembrete{
				Descricao:  "Lembrete de teste",
				Tipo:       tipo,
				Data:       "2026-08-25",
				Horario:    "14:30",
				Recorrente: false,
			}

			if err := lembrete.Validate(); err != nil {
				t.Fatalf(
					"tipo %q deveria ser válido: %v",
					tipo,
					err,
				)
			}
		})
	}
}