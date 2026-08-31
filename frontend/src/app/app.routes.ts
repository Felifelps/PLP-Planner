import { Routes } from '@angular/router';

export const routes: Routes = [
  { path: '', redirectTo: 'home', pathMatch: 'full' },
  {
    path: 'home',
    title: 'Início - PLP Planner',
    loadComponent: () =>
      import('./features/relatorios/relatorios-component').then((m) => m.RelatoriosComponent),
  },
  {
   path: 'tarefas',
   title: 'Tarefas de Hoje - PLP Planner',
   loadComponent: () =>
   import('./features/tarefas/tarefas-diarias').then((m) => m.TarefasDiarias),  },
  {
    path: 'metas',
    loadComponent: () => import('./features/metas/metas-list/metas-list').then((m) => m.MetasList),
  },
  {
    path: 'metas/nova',
    loadComponent: () => import('./features/metas/meta-form/meta-form').then((m) => m.MetaForm),
  },
  {
    path: 'metas/:id/editar',
    loadComponent: () => import('./features/metas/meta-form/meta-form').then((m) => m.MetaForm),
  },
  { path: 'relatorios', redirectTo: 'home', pathMatch: 'full' },
  { path: '**', redirectTo: 'home' },
];
