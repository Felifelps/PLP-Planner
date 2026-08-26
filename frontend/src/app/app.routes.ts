import { Routes } from '@angular/router';

export const routes: Routes = [
  { path: '', redirectTo: 'metas', pathMatch: 'full' },
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
