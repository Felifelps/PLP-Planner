import { Meta } from '../models/meta.model';
import { formatarDataLocal, paraDataInput } from './date-format.util';

export type TipoPeriodoVisualizacao = 'dia' | 'semana' | 'mes' | 'ano';

export interface IntervaloData {
  inicio: Date;
  fim: Date;
}

export function intervaloDia(referencia: Date): IntervaloData {
  const base = new Date(referencia.getFullYear(), referencia.getMonth(), referencia.getDate());
  return { inicio: base, fim: base };
}

export function intervaloSemana(referencia: Date): IntervaloData {
  const base = new Date(referencia.getFullYear(), referencia.getMonth(), referencia.getDate());
  const dia = base.getDay();
  const deslocamentoParaSegunda = dia === 0 ? -6 : 1 - dia;

  const inicio = new Date(base);
  inicio.setDate(base.getDate() + deslocamentoParaSegunda);

  const fim = new Date(inicio);
  fim.setDate(inicio.getDate() + 6);

  return { inicio, fim };
}

export function intervaloMes(referencia: Date): IntervaloData {
  const inicio = new Date(referencia.getFullYear(), referencia.getMonth(), 1);
  const fim = new Date(referencia.getFullYear(), referencia.getMonth() + 1, 0);
  return { inicio, fim };
}

export function intervaloAno(referencia: Date): IntervaloData {
  return {
    inicio: new Date(referencia.getFullYear(), 0, 1),
    fim: new Date(referencia.getFullYear(), 11, 31),
  };
}

export function calcularIntervalo(
  tipo: TipoPeriodoVisualizacao,
  referencia: Date,
): IntervaloData {
  switch (tipo) {
    case 'dia':
      return intervaloDia(referencia);
    case 'semana':
      return intervaloSemana(referencia);
    case 'mes':
      return intervaloMes(referencia);
    case 'ano':
      return intervaloAno(referencia);
  }
}

export function navegarReferencia(
  referencia: Date,
  tipo: TipoPeriodoVisualizacao,
  direcao: 1 | -1,
): Date {
  const nova = new Date(referencia);

  if (tipo === 'dia') {
    nova.setDate(nova.getDate() + direcao);
  } else if (tipo === 'semana') {
    nova.setDate(nova.getDate() + 7 * direcao);
  } else if (tipo === 'mes') {
    nova.setMonth(nova.getMonth() + direcao);
  } else {
    nova.setFullYear(nova.getFullYear() + direcao);
  }

  return nova;
}

export function metaSobrepoeIntervalo(
  meta: Pick<Meta, 'data_inicio' | 'data_fim'>,
  intervalo: IntervaloData,
): boolean {
  const metaInicio = paraDataInput(meta.data_inicio);
  const metaFim = paraDataInput(meta.data_fim);
  const intInicio = formatarDataLocal(intervalo.inicio);
  const intFim = formatarDataLocal(intervalo.fim);

  return metaInicio <= intFim && metaFim >= intInicio;
}

export function metaCobreIntervaloInteiro(
  meta: Pick<Meta, 'data_inicio' | 'data_fim'>,
  intervalo: IntervaloData,
): boolean {
  const metaInicio = paraDataInput(meta.data_inicio);
  const metaFim = paraDataInput(meta.data_fim);
  const intInicio = formatarDataLocal(intervalo.inicio);
  const intFim = formatarDataLocal(intervalo.fim);

  return metaInicio <= intInicio && metaFim >= intFim;
}