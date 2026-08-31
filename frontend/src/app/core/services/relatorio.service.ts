import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Relatorio } from '../models/relatorio.model';

@Injectable({ providedIn: 'root' })
export class RelatorioService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = `${environment.apiUrl}/relatorios`;

  gerar(dataInicio: string, dataFim: string): Observable<Relatorio> {
    const params = new HttpParams()
      .set('data_inicio', dataInicio)
      .set('data_fim', dataFim);

    return this.http.get<Relatorio>(this.baseUrl, { params });
  }
}
