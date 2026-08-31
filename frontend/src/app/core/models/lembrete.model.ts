export type TipoLembrete =
  | 'reunião'
  | 'ligação'
  | 'compra'
  | 'estudo'
  | 'exercício'
  | 'entrega';

export interface Lembrete {
  id: number;
  descricao: string;
  tipo: TipoLembrete;
  data: string;   
  horario: string;
  recorrente: boolean;
}

export type LembretePayload = Omit<Lembrete, 'id'>;
