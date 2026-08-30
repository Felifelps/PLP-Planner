import { Component, signal } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [RouterLink],
  templateUrl: './home-component.html',
  styleUrl: './home-component.css'
})
export class HomeComponent {
  protected resumo = signal({
    tarefasPendentes: 4,
    tarefasConcluidas: 6,
    metasAndamento: 2,
    lembretesHoje: 1,
    produtividade: 60 
  });
}