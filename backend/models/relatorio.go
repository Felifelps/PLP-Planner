package models

type PeriodoContagem struct {
	Rotulo string `json:"rotulo"`
	Total  int    `json:"total"`
}

type TurnoContagem struct {
	Turno string `json:"turno"`
	Total int    `json:"total"`
}

type CategoriaContagem struct {
	CategoriaID int64 `json:"categoria_id"`
	Total       int   `json:"total"`
}

type Relatorio struct {
	DataInicio   string `json:"data_inicio"`
	DataFim      string `json:"data_fim"`
	TotalMetas   int    `json:"total_metas"`
	TotalTarefas int    `json:"total_tarefas"`

	PercentualMetasCumpridas    float64 `json:"percentual_metas_cumpridas"`
	PercentualTarefasExecutadas float64 `json:"percentual_tarefas_executadas"`

	SemanaMaisProdutiva *PeriodoContagem `json:"semana_mais_produtiva,omitempty"`
	MesMaisProdutivo    *PeriodoContagem `json:"mes_mais_produtivo,omitempty"`
	TurnoMaisProdutivo  *TurnoContagem   `json:"turno_mais_produtivo,omitempty"`

	CategoriasMaisRealizadasTarefas []CategoriaContagem `json:"categorias_mais_realizadas_tarefas"`
	CategoriasMaisRealizadasMetas   []CategoriaContagem `json:"categorias_mais_realizadas_metas"`
}
