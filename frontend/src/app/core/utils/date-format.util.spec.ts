import { describe, expect, it } from 'vitest';
import { formatarDataLocal, paraDataApi, paraDataInput } from './date-format.util';

describe('date-format.util', () => {
  it('converte data de input para o formato RFC3339 esperado pelo backend', () => {
    expect(paraDataApi('2026-08-26')).toBe('2026-08-26T00:00:00Z');
  });

  it('extrai a parte de data de uma string RFC3339 vinda do backend', () => {
    expect(paraDataInput('2026-08-26T00:00:00Z')).toBe('2026-08-26');
  });

  it('faz o round-trip data-input -> api -> input sem perdas', () => {
    const original = '2026-01-05';
    expect(paraDataInput(paraDataApi(original))).toBe(original);
  });

  it('formata um Date local sem deslocar por fuso horário (evita o bug do toISOString)', () => {
    const data = new Date(2026, 0, 5); // 5 de janeiro de 2026, horário local
    expect(formatarDataLocal(data)).toBe('2026-01-05');
  });

  it('preenche dia e mês com zero à esquerda', () => {
    const data = new Date(2026, 2, 1); // 1 de março
    expect(formatarDataLocal(data)).toBe('2026-03-01');
  });
});
