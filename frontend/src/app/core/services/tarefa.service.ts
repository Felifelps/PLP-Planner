import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { StatusTarefa, Tarefa, TarefaPayload } from '../models/tarefa.model';

@Injectable({ providedIn: 'root' })
export class TarefaService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = `${environment.apiUrl}/tarefas`;

  buscarPorData(data: string): Observable<Tarefa[]> {
    return this.http.get<Tarefa[]>(`${this.baseUrl}?data=${data}`);
  }

  buscarPorId(id: number): Observable<Tarefa> {
    return this.http.get<Tarefa>(`${this.baseUrl}/${id}`);
  }

  criar(tarefa: TarefaPayload): Observable<Tarefa> {
    return this.http.post<Tarefa>(this.baseUrl, tarefa);
  }

  atualizar(id: number, tarefa: TarefaPayload): Observable<Tarefa> {
    return this.http.put<Tarefa>(`${this.baseUrl}/${id}`, tarefa);
  }

  atualizarStatus(id: number, status: StatusTarefa): Observable<void> {
    return this.http.patch<void>(`${this.baseUrl}/${id}/status`, { status });
  }

  excluir(id: number): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/${id}`);
  }
}