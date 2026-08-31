import { TitleCasePipe, UpperCasePipe } from '@angular/common';
import { Component, computed, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterLink } from '@angular/router';
import { Categoria } from '../../core/models/categoria.model';
import { Lembrete, LembretePayload, TipoLembrete } from '../../core/models/lembrete.model';
import {
  DuracaoTarefa,
  PrioridadeTarefa,
  StatusTarefa,
  Tarefa,
  TarefaPayload,
  TurnoTarefa,
} from '../../core/models/tarefa.model';
import { CategoriaService } from '../../core/services/categoria.service';
import { LembreteService } from '../../core/services/lembrete.service';
import { TarefaService } from '../../core/services/tarefa.service';

export type TipoGranularidade = '30min' | '1h' | 'turno';

export interface BlocoTempo {
  rotulo: string;
  inicio: string;
  fim: string;
  turno?: TurnoTarefa;
}

const STATUS_OPCOES: StatusTarefa[] = [
  'pendente',
  'executada',
  'parcialmente executada',
  'adiada',
  'cancelada',
];

const PRIORIDADE_OPCOES: PrioridadeTarefa[] = ['baixa', 'média', 'alta'];

const TIPOS_LEMBRETE: TipoLembrete[] = [
  'reunião',
  'ligação',
  'compra',
  'estudo',
  'exercício',
  'entrega',
];

const CORES_TIPO_LEMBRETE: Record<TipoLembrete, string> = {
  reunião: '#d0ebff',
  ligação: '#eebefa',
  compra: '#fff3bf',
  estudo: '#c3fae8',
  exercício: '#d3f9d8',
  entrega: '#ffc9c9',
};

@Component({
  selector: 'app-tarefas-diarias',
  standalone: true,
  imports: [ReactiveFormsModule, TitleCasePipe, UpperCasePipe, RouterLink],
  templateUrl: './tarefas-diarias.html',
  styleUrl: './tarefas-diarias.css',
})
export class TarefasDiarias {
  private readonly fb = inject(FormBuilder);
  private readonly tarefaService = inject(TarefaService);
  private readonly lembreteService = inject(LembreteService);
  private readonly categoriaService = inject(CategoriaService);

  protected readonly statusOpcoes = STATUS_OPCOES;
  protected readonly prioridadeOpcoes = PRIORIDADE_OPCOES;
  protected readonly tiposLembreteOpcoes = TIPOS_LEMBRETE;

  protected readonly dataReferencia = signal(new Date());
  protected readonly granularidade = signal<TipoGranularidade>('1h');
  protected readonly tarefas = signal<Tarefa[]>([]);
  protected readonly lembretes = signal<Lembrete[]>([]);
  protected readonly categorias = signal<Categoria[]>([]);
  protected readonly exibindoModalTarefa = signal(false);
  protected readonly tipoAgendamento = signal<'horario' | 'turno'>('horario');

  protected readonly categoriaPorId = computed(
    () => new Map(this.categorias().map((c) => [c.id, c]))
  );

  protected readonly dataFormatada = computed(() => {
    return new Intl.DateTimeFormat('pt-BR', {
      weekday: 'long',
      day: '2-digit',
      month: 'long',
      year: 'numeric',
    }).format(this.dataReferencia());
  });

  protected readonly dataIso = computed(() => {
    const d = this.dataReferencia();
    const ano = d.getFullYear();
    const mes = String(d.getMonth() + 1).padStart(2, '0');
    const dia = String(d.getDate()).padStart(2, '0');
    return `${ano}-${mes}-${dia}`;
  });

  protected readonly blocosGrade = computed<BlocoTempo[]>(() => {
    const modo = this.granularidade();

    if (modo === '30min') {
      const blocos: BlocoTempo[] = [];
      for (let h = 8; h <= 18; h++) {
        const hPad = h.toString().padStart(2, '0');
        const hProx = (h + 1).toString().padStart(2, '0');
        blocos.push({ rotulo: `${hPad}:00`, inicio: `${hPad}:00`, fim: `${hPad}:30` });
        blocos.push({ rotulo: `${hPad}:30`, inicio: `${hPad}:30`, fim: `${hProx}:00` });
      }
      return blocos;
    }

    if (modo === '1h') {
      return Array.from({ length: 11 }, (_, i) => {
        const h = i + 8;
        const hPad = h.toString().padStart(2, '0');
        const hProx = (h + 1).toString().padStart(2, '0');
        return { rotulo: `${hPad}:00`, inicio: `${hPad}:00`, fim: `${hProx}:00` };
      });
    }

    return [
      { rotulo: 'Manhã', inicio: '06:00', fim: '12:00', turno: 'manhã' },
      { rotulo: 'Tarde', inicio: '12:00', fim: '18:00', turno: 'tarde' },
      { rotulo: 'Noite', inicio: '18:00', fim: '23:59', turno: 'noite' },
    ];
  });

  protected readonly formTarefa = this.fb.nonNullable.group({
    descricao: ['', Validators.required],
    data: [this.dataIso(), Validators.required],
    categoriaId: [1, [Validators.required, Validators.min(1)]],
    prioridade: ['média' as PrioridadeTarefa, Validators.required],
    status: ['pendente' as StatusTarefa, Validators.required],
    horarioInicio: ['09:00'],
    duracao: ['1h' as DuracaoTarefa],
    turno: ['manhã' as TurnoTarefa],
  });

  protected readonly formLembrete = this.fb.nonNullable.group({
    descricao: ['', Validators.required],
    tipo: ['compra' as TipoLembrete, Validators.required],
    horario: ['10:00', Validators.required],
    recorrente: [false],
  });

  constructor() {
    this.carregarCategorias();
    this.carregarDadosDoDia();
  }

  private carregarCategorias(): void {
    this.categoriaService.listarTodas().subscribe({
      next: (cats) => this.categorias.set(cats ?? []),
      error: (err) => console.error('Erro ao buscar categorias do backend:', err),
    });
  }

  private carregarDadosDoDia(): void {
    const dataAtual = this.dataIso();

    this.tarefaService.buscarPorData(dataAtual).subscribe({
      next: (dados) => this.tarefas.set(dados ?? []),
      error: (err) => console.error('Erro ao buscar tarefas:', err),
    });

    this.lembreteService.buscarTodos(dataAtual, dataAtual).subscribe({
      next: (dados) => this.lembretes.set(dados ?? []),
      error: (err) => {
        console.warn('Endpoint de lembretes indisponível ou vazio:', err);
        this.lembretes.set([]);
      },
    });
  }

  protected navegarDia(delta: number): void {
    const proxima = new Date(this.dataReferencia());
    proxima.setDate(proxima.getDate() + delta);
    this.dataReferencia.set(proxima);
    this.carregarDadosDoDia();
  }

  protected irParaHoje(): void {
    this.dataReferencia.set(new Date());
    this.carregarDadosDoDia();
  }

  protected abrirModalNovaTarefa(): void {
    this.formTarefa.patchValue({ data: this.dataIso() });
    this.exibindoModalTarefa.set(true);
  }

  protected tarefasDoBloco(bloco: BlocoTempo): Tarefa[] {
    if (this.granularidade() === 'turno') {
      return this.tarefas().filter((t) => {
        if (t.turno) return t.turno === bloco.turno;
        if (t.horario_inicio) return t.horario_inicio >= bloco.inicio && t.horario_inicio < bloco.fim;
        return false;
      });
    }

    return this.tarefas().filter(
      (t) => t.horario_inicio && t.horario_inicio >= bloco.inicio && t.horario_inicio < bloco.fim
    );
  }

  protected tarefasPorTurno(): Tarefa[] {
    return this.tarefas().filter((t) => !!t.turno && this.granularidade() !== 'turno');
  }

  protected corDaCategoria(idCategoria: number): string {
    return this.categoriaPorId().get(idCategoria)?.cor ?? '#adb5bd';
  }

  protected nomeDaCategoria(idCategoria: number): string {
    return this.categoriaPorId().get(idCategoria)?.nome ?? 'Geral';
  }

  protected corDoTipoLembrete(tipo: TipoLembrete): string {
    return CORES_TIPO_LEMBRETE[tipo] ?? '#fff3bf';
  }

  protected mudarStatusTarefa(tarefa: Tarefa, novoStatus: StatusTarefa): void {
    this.tarefaService.atualizarStatus(tarefa.id, novoStatus).subscribe({
      next: () => {
        this.tarefas.update((lista) =>
          lista.map((t) => (t.id === tarefa.id ? { ...t, status: novoStatus } : t))
        );
      },
      error: (err) => console.error('Erro ao atualizar status da tarefa:', err),
    });
  }

  protected excluirLembrete(id: number): void {
    this.lembreteService.excluir(id).subscribe({
      next: () => {
        this.lembretes.update((lista) => lista.filter((l) => l.id !== id));
      },
      error: (err) => console.error('Erro ao excluir lembrete:', err),
    });
  }

  protected salvarTarefa(): void {
    if (this.formTarefa.invalid) return;
    const v = this.formTarefa.getRawValue();

    const payload: TarefaPayload = {
      descricao: v.descricao,
      categoria_id: Number(v.categoriaId),
      data: v.data,
      status: v.status,
      prioridade: v.prioridade,
      ...(this.tipoAgendamento() === 'horario'
        ? { horario_inicio: v.horarioInicio, duracao: v.duracao }
        : { turno: v.turno }),
    };

    this.tarefaService.criar(payload).subscribe({
      next: (tarefaCriada) => {
        if (tarefaCriada.data !== this.dataIso()) {
          const [ano, mes, dia] = tarefaCriada.data.split('-').map(Number);
          this.dataReferencia.set(new Date(ano, mes - 1, dia));
          this.carregarDadosDoDia();
        } else {
          this.tarefas.update((lista) => [...lista, tarefaCriada]);
        }

        this.formTarefa.reset({
          categoriaId: 1,
          prioridade: 'média',
          status: 'pendente',
          horarioInicio: '09:00',
          duracao: '1h',
          turno: 'manhã',
        });
        this.exibindoModalTarefa.set(false);
      },
      error: (err) => console.error('Erro ao salvar tarefa no backend:', err),
    });
  }

  protected adicionarLembrete(): void {
    if (this.formLembrete.invalid) return;
    const v = this.formLembrete.getRawValue();

    const payload: LembretePayload = {
      descricao: v.descricao,
      tipo: v.tipo,
      data: this.dataIso(),
      horario: v.horario,
      recorrente: v.recorrente,
    };

    this.lembreteService.criar(payload).subscribe({
      next: (lembreteCriado) => {
        this.lembretes.update((lista) => [...lista, lembreteCriado]);
        this.formLembrete.reset({
          tipo: 'compra',
          horario: '10:00',
          recorrente: false,
        });
      },
      error: (err) => console.error('Erro ao criar lembrete no backend:', err),
    });
  }
}
