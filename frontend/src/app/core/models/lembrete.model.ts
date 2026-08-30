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
  data: string;      // "YYYY-MM-DD"
  horario: string;   // "HH:MM"
  recorrente: boolean;
}

export type LembretePayload = Omit<Lembrete, 'id'>;