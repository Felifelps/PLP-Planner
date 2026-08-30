package models

import (
	"errors"
	"regexp"
	"time"
)

type StatusTarefa string

const (
	StatusExecutada             StatusTarefa = "executada"
	StatusParcialmenteExecutada StatusTarefa = "parcialmente executada"
	StatusCancelada             StatusTarefa = "cancelada"
	StatusAdiada                StatusTarefa = "adiada"
	StatusPendente              StatusTarefa = "pendente"
)

type Prioridade string

const (
	PrioridadeBaixa Prioridade = "baixa"
	PrioridadeMedia Prioridade = "média"
	PrioridadeAlta  Prioridade = "alta"
)

const (
	Duracao30Min = "30min"
	Duracao1Hora = "1h"
)

const (
	TurnoManha = "manhã"
	TurnoTarde = "tarde"
	TurnoNoite = "noite"
)

var horarioRegex = regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d$`)

type Tarefa struct {
	ID            int64        `json:"id"`
	Descricao     string       `json:"descricao"`
	CategoriaID   int64        `json:"categoria_id"`
	Data          time.Time    `json:"data"`
	HorarioInicio *string      `json:"horario_inicio,omitempty"`
	Duracao       *string      `json:"duracao,omitempty"`
	Turno         *string      `json:"turno,omitempty"`
	Status        StatusTarefa `json:"status"`
	Prioridade    Prioridade   `json:"prioridade"`
}

func (t *Tarefa) Validate() error {
	if !StatusTarefaValido(t.Status) {
		return errors.New("status inválido")
	}

	if !PrioridadeValida(t.Prioridade) {
		return errors.New("prioridade inválida")
	}

	temHorario := t.HorarioInicio != nil || t.Duracao != nil
	temTurno := t.Turno != nil

	if temHorario && temTurno {
		return errors.New("informe um horário ou um turno, não os dois")
	}

	if !temHorario && !temTurno {
		return errors.New("informe um horário ou um turno")
	}

	if temHorario {
		if t.HorarioInicio == nil || t.Duracao == nil {
			return errors.New("horário e duração devem ser informados juntos")
		}

		if !horarioRegex.MatchString(*t.HorarioInicio) {
			return errors.New("horário inválido, utilize o formato HH:MM")
		}

		if !DuracaoValida(*t.Duracao) {
			return errors.New("duração inválida")
		}
	}

	if temTurno && !TurnoValido(*t.Turno) {
		return errors.New("turno inválido")
	}

	return nil
}

func StatusTarefaValido(status StatusTarefa) bool {
	switch status {
	case StatusExecutada,
		StatusParcialmenteExecutada,
		StatusCancelada,
		StatusPendente,
		StatusAdiada:
		return true
	default:
		return false
	}
}

func PrioridadeValida(prioridade Prioridade) bool {
	switch prioridade {
	case PrioridadeBaixa,
		PrioridadeMedia,
		PrioridadeAlta:
		return true
	default:
		return false
	}
}

func DuracaoValida(duracao string) bool {
	switch duracao {
	case Duracao30Min,
		Duracao1Hora:
		return true
	default:
		return false
	}
}

func TurnoValido(turno string) bool {
	switch turno {
	case TurnoManha,
		TurnoTarde,
		TurnoNoite:
		return true
	default:
		return false
	}
}
