import { describe, expect, it } from 'vitest';
import { construirGradeMes } from './calendar.util';
import { formatarDataLocal } from './date-format.util';

describe('construirGradeMes', () => {
  it('gera semanas completas de segunda a domingo (6 semanas para agosto de 2026)', () => {
    const grade = construirGradeMes(new Date(2026, 7, 15));

    expect(grade).toHaveLength(6);
    for (const semana of grade) {
      expect(semana).toHaveLength(7);
    }

    expect(formatarDataLocal(grade[0][0].data)).toBe('2026-07-27'); // primeira segunda da grade
    expect(formatarDataLocal(grade[5][6].data)).toBe('2026-09-06'); // último domingo da grade
  });

  it('marca os dias fora do mês de referência', () => {
    const grade = construirGradeMes(new Date(2026, 7, 15));

    expect(grade[0][0].foraDoMes).toBe(true); // 27/07
    expect(grade[0][5].foraDoMes).toBe(false); // 01/08 (sábado)
    expect(grade[5][0].foraDoMes).toBe(false); // 31/08 (segunda)
    expect(grade[5][1].foraDoMes).toBe(true); // 01/09
  });

  it('gera 5 semanas quando o mês cabe exatamente nessa quantidade (fevereiro de 2026)', () => {
    const grade = construirGradeMes(new Date(2026, 1, 10));

    expect(grade).toHaveLength(5);
    expect(formatarDataLocal(grade[0][0].data)).toBe('2026-01-26');
    expect(formatarDataLocal(grade[4][6].data)).toBe('2026-03-01');
  });

  it('inclui todos os dias do mês exatamente uma vez', () => {
    const grade = construirGradeMes(new Date(2026, 7, 1));
    const diasDoMes = grade.flat().filter((dia) => !dia.foraDoMes);

    expect(diasDoMes).toHaveLength(31);
    expect(formatarDataLocal(diasDoMes[0].data)).toBe('2026-08-01');
    expect(formatarDataLocal(diasDoMes[30].data)).toBe('2026-08-31');
  });
});
