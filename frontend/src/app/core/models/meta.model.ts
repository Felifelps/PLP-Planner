export type MetaStatus = 'cumprida' | 'parcialmente cumprida' | 'não cumprida';

export type MetaPeriodo = 'semanal' | 'mensal' | 'anual';

export interface Meta {
  id: number;
  nome: string;
  descricao: string;
  categoria_id: number;
  status: MetaStatus;
  periodo: MetaPeriodo;
  data_inicio: string;
  data_fim: string;
}

export type MetaPayload = Omit<Meta, 'id'>;
