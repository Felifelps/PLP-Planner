import { Component, OnInit, inject, signal } from '@angular/core';
import { TarefaService } from '../../core/services/tarefa.service';
import { MetaService } from '../../core/services/meta.service';
import { LembreteService } from '../../core/services/lembrete.service';
import { formatarDataLocal } from '../../core/utils/date-format.util';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  private readonly tarefaService = inject(TarefaService);
  private readonly metaService = inject(MetaService);
  private readonly lembreteService = inject(LembreteService);

  tarefasPendentes = signal(0);
  tarefasConcluidas = signal(0);
  metasEmAndamento = signal(0);
  lembretesHoje = signal(0);
  produtividade = signal(0);

  ngOnInit(): void {
    this.carregarResumoDoDia();
  }

  private carregarResumoDoDia(): void {
    const hoje = formatarDataLocal(new Date());

    this.tarefaService.buscarPorData(hoje).subscribe({
      next: (tarefas) => {
        const pendentes = tarefas.filter(t => t.status === 'pendente').length;
        const concluidas = tarefas.filter(t => t.status === 'executada').length;
        
        this.tarefasPendentes.set(pendentes);
        this.tarefasConcluidas.set(concluidas);

        const total = pendentes + concluidas;
        this.produtividade.set(total > 0 ? Math.round((concluidas / total) * 100) : 0);
      },
      error: (err) => console.error('Erro ao buscar tarefas:', err)
    });

    this.lembreteService.buscarTodos(hoje, hoje).subscribe({
      next: (lembretes) => this.lembretesHoje.set(lembretes.length ?? 0),
      error: (err) => {
        console.error('Erro ao buscar lembretes:', err);
        this.lembretesHoje.set(0);
      }
    });

    this.metaService.listarTodas().subscribe({
      next: (metas) => {
        const ativas = metas.filter(m => m.status !== 'cumprida').length;
        this.metasEmAndamento.set(ativas);
      },
      error: (err) => console.error('Erro ao buscar metas:', err)
    });
  }
}