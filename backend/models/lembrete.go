package models

import "errors"

type TipoLembrete string

const (
	TipoReuniao   TipoLembrete = "reunião"
	TipoLigacao   TipoLembrete = "ligação"
	TipoCompra    TipoLembrete = "compra"
	TipoEstudo    TipoLembrete = "estudo"
	TipoExercicio TipoLembrete = "exercício"
	TipoEntrega   TipoLembrete = "entrega"
)

type Lembrete struct {
	ID         int64        `json:"id"`
	Descricao  string       `json:"descricao"`
	Tipo       TipoLembrete `json:"tipo"`
	Data       string       `json:"data"`
	Horario    string       `json:"horario"`
	Recorrente bool         `json:"recorrente"`
}

func (l *Lembrete) Validate() error {
	if !TipoLembreteValido(l.Tipo) {
		return errors.New("tipo de lembrete inválido")
	}

	return nil
}

func TipoLembreteValido(tipo TipoLembrete) bool {
	switch tipo {
	case TipoReuniao,
		TipoLigacao,
		TipoCompra,
		TipoEstudo,
		TipoExercicio,
		TipoEntrega:
		return true
	default:
		return false
	}
}