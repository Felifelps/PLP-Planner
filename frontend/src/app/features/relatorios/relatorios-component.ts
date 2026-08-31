import { CommonModule } from '@angular/common';
import { Component, OnInit, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterLink } from '@angular/router';
import { Categoria } from '../../core/models/categoria.model';
import { Lembrete } from '../../core/models/lembrete.model';
import { CategoriaContagem, Relatorio } from '../../core/models/relatorio.model';
import { CategoriaService } from '../../core/services/categoria.service';
import { LembreteService } from '../../core/services/lembrete.service';
import { MetaService } from '../../core/services/meta.service';
import { RelatorioService } from '../../core/services/relatorio.service';
import { TarefaService } from '../../core/services/tarefa.service';
import { formatarDataLocal } from '../../core/utils/date-format.util';
import {
  IntervaloData,
  TipoPeriodoVisualizacao,
  calcularIntervalo,
  navegarReferencia,
} from '../../core/utils/date-range.util';

@Component({
  selector: 'app-relatorios',
  standalone: true,
  imports: [CommonModule, RouterLink, FormsModule],
  templateUrl: './relatorios-component.html',
  styleUrl: './relatorios-component.css',
})
export class RelatoriosComponent implements OnInit {
  private readonly relatorioService = inject(RelatorioService);
  private readonly categoriaService = inject(CategoriaService);
  private readonly tarefaService = inject(TarefaService);
  private readonly metaService = inject(MetaService);
  private readonly lembreteService = inject(LembreteService);

  readonly tipoPeriodo = signal<TipoPeriodoVisualizacao>('mes');
  readonly dataReferencia = signal<Date>(new Date());
  readonly intervalo = signal<IntervaloData>(calcularIntervalo('mes', new Date()));

  readonly relatorio = signal<Relatorio | null>(null);
  readonly categorias = signal<Categoria[]>([]);
  readonly carregando = signal<boolean>(false);
  readonly erro = signal<string | null>(null);

  // Resumo do dia (bloco no topo da home)
  readonly resumoTarefasPendentes = signal(0);
  readonly resumoTarefasConcluidas = signal(0);
  readonly resumoMetasEmAndamento = signal(0);
  readonly resumoLembretesHoje = signal(0);
  readonly resumoProdutividade = signal(0);
  readonly proximosLembretes = signal<Lembrete[]>([]);

  ngOnInit(): void {
    this.carregarCategorias();
    this.carregarResumoDoDia();
    this.atualizarIntervaloECarregar();
  }

  private carregarResumoDoDia(): void {
    const hoje = formatarDataLocal(new Date());

    this.tarefaService.buscarPorData(hoje).subscribe({
      next: (tarefas) => {
        const pendentes = (tarefas ?? []).filter((t) => t.status === 'pendente').length;
        const concluidas = (tarefas ?? []).filter((t) => t.status === 'executada').length;

        this.resumoTarefasPendentes.set(pendentes);
        this.resumoTarefasConcluidas.set(concluidas);

        const total = pendentes + concluidas;
        this.resumoProdutividade.set(total > 0 ? Math.round((concluidas / total) * 100) : 0);
      },
      error: (err) => console.error('Erro ao buscar tarefas do dia:', err),
    });

    const limite = new Date();
    limite.setDate(limite.getDate() + 14);

    this.lembreteService.buscarTodos(hoje, formatarDataLocal(limite)).subscribe({
      next: (lembretes) => {
        const ordenados = [...(lembretes ?? [])].sort((a, b) =>
          a.data === b.data ? a.horario.localeCompare(b.horario) : a.data.localeCompare(b.data)
        );

        this.resumoLembretesHoje.set(ordenados.filter((l) => l.data === hoje).length);
        this.proximosLembretes.set(ordenados.slice(0, 5));
      },
      error: (err) => {
        console.error('Erro ao buscar lembretes:', err);
        this.resumoLembretesHoje.set(0);
        this.proximosLembretes.set([]);
      },
    });

    this.metaService.listarTodas().subscribe({
      next: (metas) =>
        this.resumoMetasEmAndamento.set((metas ?? []).filter((m) => m.status !== 'cumprida').length),
      error: (err) => console.error('Erro ao buscar metas:', err),
    });
  }

  selecionarPeriodo(tipo: TipoPeriodoVisualizacao): void {
    this.tipoPeriodo.set(tipo);
    this.atualizarIntervaloECarregar();
  }

  navegar(direcao: 1 | -1): void {
    const novaRef = navegarReferencia(this.dataReferencia(), this.tipoPeriodo(), direcao);
    this.dataReferencia.set(novaRef);
    this.atualizarIntervaloECarregar();
  }

  irParaHoje(): void {
    this.dataReferencia.set(new Date());
    this.atualizarIntervaloECarregar();
  }

  atualizarIntervaloECarregar(): void {
    const novoIntervalo = calcularIntervalo(this.tipoPeriodo(), this.dataReferencia());
    this.intervalo.set(novoIntervalo);
    this.carregarRelatorio();
  }

  carregarRelatorio(): void {
    const int = this.intervalo();
    const dataInicio = formatarDataLocal(int.inicio);
    const dataFim = formatarDataLocal(int.fim);

    this.carregando.set(true);
    this.erro.set(null);

    this.relatorioService.gerar(dataInicio, dataFim).subscribe({
      next: (dados) => {
        this.relatorio.set(dados);
        this.carregando.set(false);
      },
      error: (err) => {
        console.error('Erro ao buscar relatório:', err);
        this.erro.set('Não foi possível carregar o relatório para este período.');
        this.carregando.set(false);
      },
    });
  }

  private carregarCategorias(): void {
    this.categoriaService.listarTodas().subscribe({
      next: (lista) => this.categorias.set(lista),
      error: (err) => console.error('Erro ao carregar categorias:', err),
    });
  }

  rotuloIntervalo(): string {
    const { inicio, fim } = this.intervalo();
    const tipo = this.tipoPeriodo();

    if (tipo === 'semana') {
      const dInicio = inicio.toLocaleDateString('pt-BR', { day: '2-digit', month: '2-digit' });
      const dFim = fim.toLocaleDateString('pt-BR', { day: '2-digit', month: '2-digit', year: 'numeric' });
      return `${dInicio} a ${dFim}`;
    }

    if (tipo === 'mes') {
      return inicio.toLocaleDateString('pt-BR', { month: 'long', year: 'numeric' });
    }

    return `Ano de ${inicio.getFullYear()}`;
  }

  nomeDaCategoria(categoriaId: number): string {
    const cat = this.categorias().find((c) => c.id === categoriaId);
    return cat ? cat.nome : `Categoria #${categoriaId}`;
  }

  corDaCategoria(categoriaId: number): string {
    const cat = this.categorias().find((c) => c.id === categoriaId);
    return cat?.cor || '#6c757d';
  }

  metasCumpridasTotal(): number {
    const rel = this.relatorio();
    if (!rel || rel.total_metas === 0) return 0;
    return Math.round((rel.total_metas * rel.percentual_metas_cumpridas) / 100);
  }

  tarefasExecutadasTotal(): number {
    const rel = this.relatorio();
    if (!rel || rel.total_tarefas === 0) return 0;
    return Math.round((rel.total_tarefas * rel.percentual_tarefas_executadas) / 100);
  }

  formatarRotuloSemana(rotulo?: string): string {
    if (!rotulo) return 'Nenhuma tarefa executada';
    const partes = rotulo.split('-W');
    if (partes.length === 2) {
      return `Semana ${partes[1]} (${partes[0]})`;
    }
    return rotulo;
  }

  formatarRotuloMes(rotulo?: string): string {
    if (!rotulo) return 'Nenhuma tarefa executada';
    const partes = rotulo.split('-');
    if (partes.length === 2) {
      const ano = parseInt(partes[0], 10);
      const mes = parseInt(partes[1], 10) - 1;
      const data = new Date(ano, mes, 1);
      return data.toLocaleDateString('pt-BR', { month: 'long', year: 'numeric' });
    }
    return rotulo;
  }

  iconeTurno(turno?: string): string {
    switch (turno?.toLowerCase()) {
      case 'manhã':
        return '🌅';
      case 'tarde':
        return '☀️';
      case 'noite':
        return '🌙';
      default:
        return '⏱️';
    }
  }

  rotuloDataLembrete(iso: string): string {
    const hoje = formatarDataLocal(new Date());
    const amanha = new Date();
    amanha.setDate(amanha.getDate() + 1);

    if (iso === hoje) return 'Hoje';
    if (iso === formatarDataLocal(amanha)) return 'Amanhã';

    const [ano, mes, dia] = iso.split('-').map(Number);
    return new Intl.DateTimeFormat('pt-BR', {
      weekday: 'short',
      day: '2-digit',
      month: '2-digit',
    })
      .format(new Date(ano, mes - 1, dia))
      .replaceAll('.', '');
  }

  calcularPorcentagemCategoria(item: CategoriaContagem, lista: CategoriaContagem[]): number {
    const totalSoma = lista.reduce((acc, curr) => acc + curr.total, 0);
    if (totalSoma === 0) return 0;
    return Math.round((item.total / totalSoma) * 100);
  }
}
