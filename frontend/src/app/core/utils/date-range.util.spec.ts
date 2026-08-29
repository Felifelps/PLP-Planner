import { describe, expect, it } from 'vitest';
import {
  calcularIntervalo,
  intervaloAno,
  intervaloMes,
  intervaloSemana,
  metaCobreIntervaloInteiro,
  metaSobrepoeIntervalo,
  navegarReferencia,
} from './date-range.util';
import { formatarDataLocal } from './date-format.util';

describe('intervaloSemana', () => {
  it('encontra a segunda-feira quando a referência é uma quarta-feira', () => {
    const quarta = new Date(2026, 7, 26); // 26/08/2026 é uma quarta-feira
    const { inicio, fim } = intervaloSemana(quarta);
    expect(formatarDataLocal(inicio)).toBe('2026-08-24');
    expect(formatarDataLocal(fim)).toBe('2026-08-30');
  });

  it('volta 6 dias quando a referência cai num domingo', () => {
    const domingo = new Date(2026, 7, 30); // 30/08/2026 é domingo
    const { inicio, fim } = intervaloSemana(domingo);
    expect(formatarDataLocal(inicio)).toBe('2026-08-24');
    expect(formatarDataLocal(fim)).toBe('2026-08-30');
  });
});

describe('intervaloMes', () => {
  it('calcula o intervalo correto em fevereiro de ano bissexto', () => {
    const referencia = new Date(2024, 1, 10);
    const { inicio, fim } = intervaloMes(referencia);
    expect(formatarDataLocal(inicio)).toBe('2024-02-01');
    expect(formatarDataLocal(fim)).toBe('2024-02-29');
  });

  it('calcula o intervalo correto em fevereiro de ano não bissexto', () => {
    const referencia = new Date(2026, 1, 10);
    const { inicio, fim } = intervaloMes(referencia);
    expect(formatarDataLocal(inicio)).toBe('2026-02-01');
    expect(formatarDataLocal(fim)).toBe('2026-02-28');
  });
});

describe('intervaloAno', () => {
  it('vai de 1º de janeiro a 31 de dezembro', () => {
    const { inicio, fim } = intervaloAno(new Date(2026, 5, 15));
    expect(formatarDataLocal(inicio)).toBe('2026-01-01');
    expect(formatarDataLocal(fim)).toBe('2026-12-31');
  });
});

describe('calcularIntervalo', () => {
  it('despacha para a função correta de acordo com o tipo', () => {
    const referencia = new Date(2026, 7, 26);
    expect(calcularIntervalo('semana', referencia)).toEqual(intervaloSemana(referencia));
    expect(calcularIntervalo('mes', referencia)).toEqual(intervaloMes(referencia));
    expect(calcularIntervalo('ano', referencia)).toEqual(intervaloAno(referencia));
  });
});

describe('navegarReferencia', () => {
  it('avança 7 dias para o tipo semana', () => {
    const referencia = new Date(2026, 7, 26);
    const nova = navegarReferencia(referencia, 'semana', 1);
    expect(formatarDataLocal(nova)).toBe('2026-09-02');
  });

  it('volta um mês para o tipo mes', () => {
    const referencia = new Date(2026, 7, 26);
    const nova = navegarReferencia(referencia, 'mes', -1);
    expect(formatarDataLocal(nova)).toBe('2026-07-26');
  });

  it('vira o ano corretamente ao avançar em dezembro', () => {
    const referencia = new Date(2026, 11, 15);
    const nova = navegarReferencia(referencia, 'ano', 1);
    expect(formatarDataLocal(nova)).toBe('2027-12-15');
  });
});

describe('metaSobrepoeIntervalo', () => {
  const intervalo = { inicio: new Date(2026, 7, 24), fim: new Date(2026, 7, 30) };

  it('retorna false para meta totalmente fora do intervalo', () => {
    const meta = { data_inicio: '2026-09-01T00:00:00Z', data_fim: '2026-09-05T00:00:00Z' };
    expect(metaSobrepoeIntervalo(meta, intervalo)).toBe(false);
  });

  it('retorna true para meta parcialmente sobreposta no início do intervalo', () => {
    const meta = { data_inicio: '2026-08-20T00:00:00Z', data_fim: '2026-08-25T00:00:00Z' };
    expect(metaSobrepoeIntervalo(meta, intervalo)).toBe(true);
  });

  it('retorna true para meta parcialmente sobreposta no fim do intervalo', () => {
    const meta = { data_inicio: '2026-08-29T00:00:00Z', data_fim: '2026-09-10T00:00:00Z' };
    expect(metaSobrepoeIntervalo(meta, intervalo)).toBe(true);
  });

  it('retorna true para meta anual que contém o intervalo inteiro', () => {
    const meta = { data_inicio: '2026-01-01T00:00:00Z', data_fim: '2026-12-31T00:00:00Z' };
    expect(metaSobrepoeIntervalo(meta, intervalo)).toBe(true);
  });
});

describe('metaCobreIntervaloInteiro', () => {
  const intervalo = { inicio: new Date(2026, 7, 24), fim: new Date(2026, 7, 30) };

  it('retorna true para meta anual que cobre a semana inteira', () => {
    const meta = { data_inicio: '2026-01-01T00:00:00Z', data_fim: '2026-12-31T00:00:00Z' };
    expect(metaCobreIntervaloInteiro(meta, intervalo)).toBe(true);
  });

  it('retorna true quando a meta cobre exatamente o intervalo', () => {
    const meta = { data_inicio: '2026-08-24T00:00:00Z', data_fim: '2026-08-30T00:00:00Z' };
    expect(metaCobreIntervaloInteiro(meta, intervalo)).toBe(true);
  });

  it('retorna false para meta que só cobre parte do intervalo', () => {
    const meta = { data_inicio: '2026-08-24T00:00:00Z', data_fim: '2026-08-27T00:00:00Z' };
    expect(metaCobreIntervaloInteiro(meta, intervalo)).toBe(false);
  });

  it('retorna false para meta que começa depois do intervalo iniciar', () => {
    const meta = { data_inicio: '2026-08-25T00:00:00Z', data_fim: '2026-09-30T00:00:00Z' };
    expect(metaCobreIntervaloInteiro(meta, intervalo)).toBe(false);
  });
});
