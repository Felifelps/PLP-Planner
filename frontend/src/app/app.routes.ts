import { Routes } from '@angular/router';

export const routes: Routes = [
  { path: '', redirectTo: 'home', pathMatch: 'full' },
  { 
    path: 'home', 
    loadComponent: () => import('./features/home/home-component').then(m => m.HomeComponent) 
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
];
