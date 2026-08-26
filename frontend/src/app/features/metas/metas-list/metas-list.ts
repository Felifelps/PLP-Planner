import { Component, computed, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { CategoriaService } from '../../../core/services/categoria.service';
import { MetaService } from '../../../core/services/meta.service';
import { Categoria } from '../../../core/models/categoria.model';
import { Meta, MetaStatus } from '../../../core/models/meta.model';
import {
  TipoPeriodoVisualizacao,
  calcularIntervalo,
  metaSobrepoeIntervalo,
  navegarReferencia,
} from '../../../core/utils/date-range.util';

const STATUS_OPCOES: MetaStatus[] = ['não cumprida', 'parcialmente cumprida', 'cumprida'];

@Component({
  selector: 'app-metas-list',
  imports: [RouterLink],
  templateUrl: './metas-list.html',
  styleUrl: './metas-list.css',
})
export class MetasList {
  private readonly metaService = inject(MetaService);
  private readonly categoriaService = inject(CategoriaService);

  protected readonly statusOpcoes = STATUS_OPCOES;

  protected readonly todasMetas = signal<Meta[]>([]);
  protected readonly categorias = signal<Categoria[]>([]);
  protected readonly tipoPeriodo = signal<TipoPeriodoVisualizacao>('semana');
  protected readonly referencia = signal(new Date());
  protected readonly carregando = signal(false);
  protected readonly erro = signal<string | null>(null);

  protected readonly intervalo = computed(() =>
    calcularIntervalo(this.tipoPeriodo(), this.referencia()),
  );

  protected readonly metasVisiveis = computed(() => {
    const intervalo = this.intervalo();
    return this.todasMetas()
      .filter((meta) => metaSobrepoeIntervalo(meta, intervalo))
      .sort((a, b) => a.data_inicio.localeCompare(b.data_inicio));
  });

  protected readonly categoriaPorId = computed(
    () => new Map(this.categorias().map((categoria) => [categoria.id, categoria])),
  );

  protected readonly rotuloIntervalo = computed(() => {
    const { inicio, fim } = this.intervalo();
    const formatador = new Intl.DateTimeFormat('pt-BR', { day: '2-digit', month: 'short', year: 'numeric' });
    return `${formatador.format(inicio)} – ${formatador.format(fim)}`;
  });

  constructor() {
    this.carregarDados();
  }

  private carregarDados(): void {
    this.carregando.set(true);
    this.erro.set(null);

    this.categoriaService.listarTodas().subscribe({
      next: (categorias) => this.categorias.set(categorias),
      error: () => this.erro.set('Não foi possível carregar as categorias.'),
    });

    this.metaService.listarTodas().subscribe({
      next: (metas) => {
        this.todasMetas.set(metas);
        this.carregando.set(false);
      },
      error: () => {
        this.erro.set('Não foi possível carregar as metas.');
        this.carregando.set(false);
      },
    });
  }

  protected selecionarPeriodo(tipo: TipoPeriodoVisualizacao): void {
    this.tipoPeriodo.set(tipo);
  }

  protected navegar(direcao: 1 | -1): void {
    this.referencia.update((referencia) => navegarReferencia(referencia, this.tipoPeriodo(), direcao));
  }

  protected corDaCategoria(meta: Meta): string {
    return this.categoriaPorId().get(meta.categoria_id)?.cor ?? '#adb5bd';
  }

  protected nomeDaCategoria(meta: Meta): string {
    return this.categoriaPorId().get(meta.categoria_id)?.nome ?? 'Sem categoria';
  }

  protected mudarStatus(meta: Meta, status: MetaStatus): void {
    this.metaService.atualizarStatus(meta.id, status).subscribe({
      next: () => {
        this.todasMetas.update((metas) =>
          metas.map((m) => (m.id === meta.id ? { ...m, status } : m)),
        );
      },
      error: () => this.erro.set('Não foi possível atualizar o status da meta.'),
    });
  }

  protected excluir(meta: Meta): void {
    if (!confirm(`Excluir a meta "${meta.nome}"?`)) {
      return;
    }

    this.metaService.excluir(meta.id).subscribe({
      next: () => {
        this.todasMetas.update((metas) => metas.filter((m) => m.id !== meta.id));
      },
      error: () => this.erro.set('Não foi possível excluir a meta.'),
    });
  }
}
