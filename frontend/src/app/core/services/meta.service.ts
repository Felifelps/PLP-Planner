import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Meta, MetaPayload, MetaStatus } from '../models/meta.model';

@Injectable({ providedIn: 'root' })
export class MetaService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = `${environment.apiUrl}/metas`;

  listarTodas(): Observable<Meta[]> {
    return this.http.get<Meta[]>(this.baseUrl);
  }

  buscarPorId(id: number): Observable<Meta> {
    return this.http.get<Meta>(`${this.baseUrl}/${id}`);
  }

  criar(meta: MetaPayload): Observable<Meta> {
    return this.http.post<Meta>(this.baseUrl, meta);
  }

  atualizar(id: number, meta: MetaPayload): Observable<Meta> {
    return this.http.put<Meta>(`${this.baseUrl}/${id}`, meta);
  }

  atualizarStatus(id: number, status: MetaStatus): Observable<void> {
    return this.http.patch<void>(`${this.baseUrl}/${id}/status`, { status });
  }

  excluir(id: number): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/${id}`);
  }
}
