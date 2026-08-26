import { intervaloSemana } from './date-range.util';

export interface DiaCalendario {
  data: Date;
  foraDoMes: boolean;
}

export const DIAS_SEMANA_ABREV = ['Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb', 'Dom'];

export function construirGradeMes(referencia: Date): DiaCalendario[][] {
  const mes = referencia.getMonth();
  const primeiroDia = new Date(referencia.getFullYear(), mes, 1);
  const ultimoDia = new Date(referencia.getFullYear(), mes + 1, 0);

  const inicioGrade = intervaloSemana(primeiroDia).inicio;
  const fimGrade = intervaloSemana(ultimoDia).fim;

  const dias: DiaCalendario[] = [];
  const cursor = new Date(inicioGrade);

  while (cursor <= fimGrade) {
    dias.push({ data: new Date(cursor), foraDoMes: cursor.getMonth() !== mes });
    cursor.setDate(cursor.getDate() + 1);
  }

  const semanas: DiaCalendario[][] = [];
  for (let i = 0; i < dias.length; i += 7) {
    semanas.push(dias.slice(i, i + 7));
  }

  return semanas;
}
