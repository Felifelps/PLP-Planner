export type StatusTarefa = 'pendente' | 'executada' | 'parcialmente executada' | 'cancelada' | 'adiada';
export type PrioridadeTarefa = 'baixa' | 'média' | 'alta';
export type DuracaoTarefa = '30min' | '1h';
export type TurnoTarefa = 'manhã' | 'tarde' | 'noite';


export interface Tarefa {
  id: number;
  descricao: string;
  categoria_id: number;
  data: string;
  horario_inicio?: string;
  duracao?: DuracaoTarefa;
  turno?: TurnoTarefa;
  status: StatusTarefa;
  prioridade: PrioridadeTarefa;
}

export type TarefaPayload = Omit<Tarefa, 'id'>;
