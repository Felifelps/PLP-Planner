import { Component, effect, inject, input, signal } from '@angular/core';
import {
  AbstractControl,
  FormBuilder,
  ReactiveFormsModule,
  ValidationErrors,
  Validators,
} from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { CategoriaService } from '../../../core/services/categoria.service';
import { MetaService } from '../../../core/services/meta.service';
import { Categoria } from '../../../core/models/categoria.model';
import { MetaPayload, MetaPeriodo, MetaStatus } from '../../../core/models/meta.model';
import { paraDataApi, paraDataInput } from '../../../core/utils/date-format.util';

function intervaloValido(grupo: AbstractControl): ValidationErrors | null {
  const dataInicio = grupo.get('dataInicio')?.value;
  const dataFim = grupo.get('dataFim')?.value;

  if (dataInicio && dataFim && dataInicio > dataFim) {
    return { intervaloInvalido: true };
  }

  return null;
}

@Component({
  selector: 'app-meta-form',
  imports: [ReactiveFormsModule, RouterLink],
  templateUrl: './meta-form.html',
  styleUrl: './meta-form.css',
})
export class MetaForm {
  readonly id = input<string>();

  private readonly fb = inject(FormBuilder);
  private readonly router = inject(Router);
  private readonly metaService = inject(MetaService);
  private readonly categoriaService = inject(CategoriaService);

  protected readonly categorias = signal<Categoria[]>([]);
  protected readonly salvando = signal(false);
  protected readonly erro = signal<string | null>(null);

  protected readonly form = this.fb.nonNullable.group(
    {
      nome: ['', Validators.required],
      descricao: [''],
      categoriaId: [0, [Validators.required, Validators.min(1)]],
      periodo: ['semanal' as MetaPeriodo, Validators.required],
      status: ['não cumprida' as MetaStatus, Validators.required],
      dataInicio: ['', Validators.required],
      dataFim: ['', Validators.required],
    },
    { validators: intervaloValido },
  );

  protected readonly modoEdicao = signal(false);

  constructor() {
    this.categoriaService.listarTodas().subscribe({
      next: (categorias) => this.categorias.set(categorias),
      error: () => this.erro.set('Não foi possível carregar as categorias.'),
    });

    effect(() => {
      const id = this.id();
      if (id) {
        this.modoEdicao.set(true);
        this.carregarParaEdicao(Number(id));
      }
    });
  }

  private carregarParaEdicao(id: number): void {
    this.metaService.buscarPorId(id).subscribe({
      next: (meta) => {
        this.form.setValue({
          nome: meta.nome,
          descricao: meta.descricao,
          categoriaId: meta.categoria_id,
          periodo: meta.periodo,
          status: meta.status,
          dataInicio: paraDataInput(meta.data_inicio),
          dataFim: paraDataInput(meta.data_fim),
        });
      },
      error: () => this.erro.set('Não foi possível carregar a meta.'),
    });
  }

  protected salvar(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const valores = this.form.getRawValue();

    const payload: MetaPayload = {
      nome: valores.nome,
      descricao: valores.descricao,
      categoria_id: valores.categoriaId,
      periodo: valores.periodo,
      status: valores.status,
      data_inicio: paraDataApi(valores.dataInicio),
      data_fim: paraDataApi(valores.dataFim),
    };

    this.salvando.set(true);
    this.erro.set(null);

    const idAtual = this.id();
    const requisicao = idAtual
      ? this.metaService.atualizar(Number(idAtual), payload)
      : this.metaService.criar(payload);

    requisicao.subscribe({
      next: () => this.router.navigate(['/metas']),
      error: (erro) => {
        this.salvando.set(false);
        this.erro.set(typeof erro?.error === 'string' ? erro.error : 'Não foi possível salvar a meta.');
      },
    });
  }
}
