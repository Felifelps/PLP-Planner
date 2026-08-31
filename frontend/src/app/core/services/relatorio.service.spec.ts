import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import { environment } from '../../../environments/environment';
import { Relatorio } from '../models/relatorio.model';
import { RelatorioService } from './relatorio.service';

describe('RelatorioService', () => {
  let service: RelatorioService;
  let httpMock: HttpTestingController;
  const baseUrl = `${environment.apiUrl}/relatorios`;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideHttpClient(), provideHttpClientTesting()],
    });

    service = TestBed.inject(RelatorioService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('deve gerar relatorio com parametros de dataInicio e dataFim', () => {
    const mockRelatorio: Relatorio = {
      data_inicio: '2026-09-01',
      data_fim: '2026-09-30',
      total_metas: 10,
      total_tarefas: 20,
      percentual_metas_cumpridas: 80,
      percentual_tarefas_executadas: 75,
      semana_mais_produtiva: { rotulo: '2026-W36', total: 10 },
      mes_mais_produtivo: { rotulo: '2026-09', total: 20 },
      turno_mais_produtivo: { turno: 'manhã', total: 12 },
      categorias_mais_realizadas_tarefas: [{ categoria_id: 1, total: 10 }],
      categorias_mais_realizadas_metas: [{ categoria_id: 1, total: 5 }],
    };

    service.gerar('2026-09-01', '2026-09-30').subscribe((relatorio) => {
      expect(relatorio).toEqual(mockRelatorio);
    });

    const req = httpMock.expectOne(`${baseUrl}?data_inicio=2026-09-01&data_fim=2026-09-30`);
    expect(req.request.method).toBe('GET');
    req.flush(mockRelatorio);
  });
});
