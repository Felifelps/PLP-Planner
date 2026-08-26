package models

import (
	"testing"
	"time"
)

func ptr(s string) *string {
	return &s
}

func TestTarefaValidate(t *testing.T) {
	data := time.Date(2026, time.August, 19, 0, 0, 0, 0, time.UTC)

	horario := ptr("14:30")
	duracao := ptr(Duracao30Min)
	turno := ptr(TurnoManha)

	testes := []struct {
		nome   string
		tarefa Tarefa
		erro   bool
	}{
		{
			nome: "tarefa válida com horário",
			tarefa: Tarefa{
				Data:          data,
				HorarioInicio: horario,
				Duracao:       duracao,
				Status:        StatusExecutada,
				Prioridade:    PrioridadeAlta,
			},
			erro: false,
		},
		{
			nome: "tarefa válida com turno",
			tarefa: Tarefa{
				Data:       data,
				Turno:      turno,
				Status:     StatusExecutada,
				Prioridade: PrioridadeAlta,
			},
			erro: false,
		},
		{
			nome: "horário e turno informados juntos",
			tarefa: Tarefa{
				Data:          data,
				HorarioInicio: horario,
				Duracao:       duracao,
				Turno:         turno,
				Status:        StatusExecutada,
				Prioridade:    PrioridadeAlta,
			},
			erro: true,
		},
		{
			nome: "nem horário nem turno informados",
			tarefa: Tarefa{
				Data:       data,
				Status:     StatusExecutada,
				Prioridade: PrioridadeAlta,
			},
			erro: true,
		},
		{
			nome: "horário sem duração",
			tarefa: Tarefa{
				Data:          data,
				HorarioInicio: horario,
				Status:        StatusExecutada,
				Prioridade:    PrioridadeAlta,
			},
			erro: true,
		},
		{
			nome: "horário em formato inválido",
			tarefa: Tarefa{
				Data:          data,
				HorarioInicio: ptr("25:99"),
				Duracao:       duracao,
				Status:        StatusExecutada,
				Prioridade:    PrioridadeAlta,
			},
			erro: true,
		},
		{
			nome: "duração inválida",
			tarefa: Tarefa{
				Data:          data,
				HorarioInicio: horario,
				Duracao:       ptr("2h"),
				Status:        StatusExecutada,
				Prioridade:    PrioridadeAlta,
			},
			erro: true,
		},
		{
			nome: "turno inválido",
			tarefa: Tarefa{
				Data:       data,
				Turno:      ptr("madrugada"),
				Status:     StatusExecutada,
				Prioridade: PrioridadeAlta,
			},
			erro: true,
		},
		{
			nome: "status inválido",
			tarefa: Tarefa{
				Data:       data,
				Turno:      turno,
				Status:     StatusTarefa("invalido"),
				Prioridade: PrioridadeAlta,
			},
			erro: true,
		},
		{
			nome: "prioridade inválida",
			tarefa: Tarefa{
				Data:       data,
				Turno:      turno,
				Status:     StatusExecutada,
				Prioridade: Prioridade("invalida"),
			},
			erro: true,
		},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			err := tt.tarefa.Validate()

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

func TestStatusTarefaValidos(t *testing.T) {
	statuses := []StatusTarefa{
		StatusExecutada,
		StatusParcialmenteExecutada,
		StatusCancelada,
		StatusAdiada,
	}

	for _, status := range statuses {
		t.Run(string(status), func(t *testing.T) {
			if !StatusTarefaValido(status) {
				t.Fatalf("status %q deveria ser válido", status)
			}
		})
	}
}

func TestPrioridadesValidas(t *testing.T) {
	prioridades := []Prioridade{
		PrioridadeBaixa,
		PrioridadeMedia,
		PrioridadeAlta,
	}

	for _, prioridade := range prioridades {
		t.Run(string(prioridade), func(t *testing.T) {
			if !PrioridadeValida(prioridade) {
				t.Fatalf("prioridade %q deveria ser válida", prioridade)
			}
		})
	}
}
