import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Categoria, CategoriaPayload } from '../models/categoria.model';

@Injectable({ providedIn: 'root' })
export class CategoriaService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = `${environment.apiUrl}/categorias`;

  listarTodas(): Observable<Categoria[]> {
    return this.http.get<Categoria[]>(this.baseUrl);
  }

  criar(categoria: CategoriaPayload): Observable<Categoria> {
    return this.http.post<Categoria>(this.baseUrl, categoria);
  }

  atualizar(id: number, categoria: CategoriaPayload): Observable<Categoria> {
    return this.http.put<Categoria>(`${this.baseUrl}/${id}`, categoria);
  }

  excluir(id: number): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/${id}`);
  }
}
