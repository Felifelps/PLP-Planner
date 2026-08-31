import { Component, computed, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { HttpErrorResponse } from '@angular/common/http';
import { Categoria, CategoriaPayload } from '../../core/models/categoria.model';
import { CategoriaService } from '../../core/services/categoria.service';

const COR_HEX = /^#[0-9A-Fa-f]{6}$/;

@Component({
  selector: 'app-categorias-page',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './categorias-page.html',
  styleUrl: './categorias-page.css',
})
export class CategoriasPage {
  private readonly fb = inject(FormBuilder);
  private readonly categoriaService = inject(CategoriaService);

  protected readonly categorias = signal<Categoria[]>([]);
  protected readonly emEdicao = signal<Categoria | null>(null);
  protected readonly erro = signal<string | null>(null);
  protected readonly carregando = signal(false);

  protected readonly tituloForm = computed(() =>
    this.emEdicao() ? 'Editar categoria' : 'Nova categoria'
  );

  protected readonly form = this.fb.nonNullable.group({
    nome: ['', [Validators.required, Validators.maxLength(100)]],
    cor: ['#4C6EF5', [Validators.required, Validators.pattern(COR_HEX)]],
  });

  constructor() {
    this.carregar();
  }

  private carregar(): void {
    this.carregando.set(true);
    this.categoriaService.listarTodas().subscribe({
      next: (lista) => {
        this.categorias.set([...(lista ?? [])].sort((a, b) => a.nome.localeCompare(b.nome)));
        this.carregando.set(false);
      },
      error: () => {
        this.erro.set('Não foi possível carregar as categorias.');
        this.carregando.set(false);
      },
    });
  }

  protected iniciarEdicao(categoria: Categoria): void {
    this.emEdicao.set(categoria);
    this.erro.set(null);
    this.form.setValue({ nome: categoria.nome, cor: categoria.cor });
  }

  protected cancelarEdicao(): void {
    this.emEdicao.set(null);
    this.erro.set(null);
    this.form.reset({ nome: '', cor: '#4C6EF5' });
  }

  protected salvar(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.erro.set(null);
    const payload: CategoriaPayload = this.form.getRawValue();
    const alvo = this.emEdicao();

    const requisicao = alvo
      ? this.categoriaService.atualizar(alvo.id, payload)
      : this.categoriaService.criar(payload);

    requisicao.subscribe({
      next: () => {
        this.cancelarEdicao();
        this.carregar();
      },
      error: (err: HttpErrorResponse) => {
        this.erro.set(
          err.status === 409
            ? 'Já existe uma categoria com esse nome.'
            : 'Não foi possível salvar a categoria.'
        );
      },
    });
  }

  protected excluir(categoria: Categoria): void {
    if (!confirm(`Excluir a categoria "${categoria.nome}"?`)) {
      return;
    }

    this.erro.set(null);
    this.categoriaService.excluir(categoria.id).subscribe({
      next: () => {
        if (this.emEdicao()?.id === categoria.id) {
          this.cancelarEdicao();
        }
        this.carregar();
      },
      error: (err: HttpErrorResponse) => {
        this.erro.set(
          err.status === 409
            ? 'Esta categoria está em uso por tarefas ou metas e não pode ser excluída.'
            : 'Não foi possível excluir a categoria.'
        );
      },
    });
  }
}
