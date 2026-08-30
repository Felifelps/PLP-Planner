import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Lembrete, LembretePayload } from '../models/lembrete.model';

@Injectable({ providedIn: 'root' })
export class LembreteService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = `${environment.apiUrl}/lembretes`;

  buscarTodos(dataInicio?: string, dataFim?: string): Observable<Lembrete[]> {
    let params = new HttpParams();
    if (dataInicio) params = params.set('data_inicio', dataInicio);
    if (dataFim) params = params.set('data_fim', dataFim);

    return this.http.get<Lembrete[]>(this.baseUrl, { params });
  }

  buscarPorId(id: number): Observable<Lembrete> {
    return this.http.get<Lembrete>(`${this.baseUrl}/${id}`);
  }

  criar(lembrete: LembretePayload): Observable<Lembrete> {
    return this.http.post<Lembrete>(this.baseUrl, lembrete);
  }

  atualizar(id: number, lembrete: LembretePayload): Observable<Lembrete> {
    return this.http.put<Lembrete>(`${this.baseUrl}/${id}`, lembrete);
  }

  excluir(id: number): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/${id}`);
  }
}