import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import { environment } from '../../../environments/environment';
import { MetaPayload } from '../models/meta.model';
import { MetaService } from './meta.service';

describe('MetaService', () => {
  let service: MetaService;
  let httpMock: HttpTestingController;
  const baseUrl = `${environment.apiUrl}/metas`;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideHttpClient(), provideHttpClientTesting()],
    });

    service = TestBed.inject(MetaService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('lista todas as metas sem query params', () => {
    service.listarTodas().subscribe();

    const req = httpMock.expectOne(baseUrl);
    expect(req.request.method).toBe('GET');
    req.flush([]);
  });

  it('busca uma meta por id', () => {
    service.buscarPorId(7).subscribe();

    const req = httpMock.expectOne(`${baseUrl}/7`);
    expect(req.request.method).toBe('GET');
    req.flush({});
  });

  it('cria uma meta enviando o payload no corpo', () => {
    const payload: MetaPayload = {
      nome: 'Ler um livro',
      descricao: 'Terminar o livro do mês',
      categoria_id: 2,
      status: 'não cumprida',
      periodo: 'mensal',
      data_inicio: '2026-08-01T00:00:00Z',
      data_fim: '2026-08-31T00:00:00Z',
    };

    service.criar(payload).subscribe();

    const req = httpMock.expectOne(baseUrl);
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual(payload);
    req.flush({ id: 1, ...payload });
  });

  it('atualiza uma meta existente', () => {
    const payload: MetaPayload = {
      nome: 'Ler um livro',
      descricao: '',
      categoria_id: 2,
      status: 'cumprida',
      periodo: 'mensal',
      data_inicio: '2026-08-01T00:00:00Z',
      data_fim: '2026-08-31T00:00:00Z',
    };

    service.atualizar(3, payload).subscribe();

    const req = httpMock.expectOne(`${baseUrl}/3`);
    expect(req.request.method).toBe('PUT');
    expect(req.request.body).toEqual(payload);
    req.flush({ id: 3, ...payload });
  });

  it('atualiza apenas o status pelo endpoint dedicado', () => {
    service.atualizarStatus(3, 'cumprida').subscribe();

    const req = httpMock.expectOne(`${baseUrl}/3/status`);
    expect(req.request.method).toBe('PATCH');
    expect(req.request.body).toEqual({ status: 'cumprida' });
    req.flush(null);
  });

  it('exclui uma meta', () => {
    service.excluir(3).subscribe();

    const req = httpMock.expectOne(`${baseUrl}/3`);
    expect(req.request.method).toBe('DELETE');
    req.flush(null);
  });
});
