export interface PeriodoContagem {
  rotulo: string;
  total: number;
}

export interface TurnoContagem {
  turno: string;
  total: number;
}

export interface CategoriaContagem {
  categoria_id: number;
  total: number;
}

export interface Relatorio {
  data_inicio: string;
  data_fim: string;
  total_metas: number;
  total_tarefas: number;
  percentual_metas_cumpridas: number;
  percentual_tarefas_executadas: number;
  semana_mais_produtiva?: PeriodoContagem;
  mes_mais_produtivo?: PeriodoContagem;
  turno_mais_produtivo?: TurnoContagem;
  categorias_mais_realizadas_tarefas: CategoriaContagem[];
  categorias_mais_realizadas_metas: CategoriaContagem[];
}
