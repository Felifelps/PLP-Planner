import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { of, throwError } from 'rxjs';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { CategoriaService } from '../../core/services/categoria.service';
import { RelatorioService } from '../../core/services/relatorio.service';
import { RelatoriosComponent } from './relatorios-component';

describe('RelatoriosComponent', () => {
  let component: RelatoriosComponent;
  let fixture: ComponentFixture<RelatoriosComponent>;
  let relatorioService: RelatorioService;
  let categoriaService: CategoriaService;

  const mockRelatorio = {
    data_inicio: '2026-09-01',
    data_fim: '2026-09-30',
    total_metas: 10,
    total_tarefas: 20,
    percentual_metas_cumpridas: 80,
    percentual_tarefas_executadas: 75,
    semana_mais_produtiva: { rotulo: '2026-W36', total: 10 },
    mes_mais_produtivo: { rotulo: '2026-09', total: 20 },
    turno_mais_produtivo: { turno: 'manhã', total: 12 },
    categorias_mais_realizadas_tarefas: [{ categoria_id: 1, total: 10 }],
    categorias_mais_realizadas_metas: [{ categoria_id: 1, total: 5 }],
  };

  const mockCategorias = [
    { id: 1, nome: 'Faculdade', cor: '#3b82f6' },
    { id: 2, nome: 'Trabalho', cor: '#10b981' },
  ];

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RelatoriosComponent],
      providers: [
        provideHttpClient(),
        provideHttpClientTesting(),
        provideRouter([]),
      ],
    }).compileComponents();

    relatorioService = TestBed.inject(RelatorioService);
    categoriaService = TestBed.inject(CategoriaService);

    vi.spyOn(relatorioService, 'gerar').mockReturnValue(of(mockRelatorio));
    vi.spyOn(categoriaService, 'listarTodas').mockReturnValue(of(mockCategorias));

    fixture = TestBed.createComponent(RelatoriosComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('deve criar o componente', () => {
    expect(component).toBeTruthy();
  });

  it('deve carregar o relatorio e categorias na inicializacao', () => {
    expect(component.relatorio()).toEqual(mockRelatorio);
    expect(component.categorias()).toEqual(mockCategorias);
    expect(component.carregando()).toBe(false);
  });

  it('deve calcular corretamente o total de metas cumpridas e tarefas executadas', () => {
    expect(component.metasCumpridasTotal()).toBe(8); 
    expect(component.tarefasExecutadasTotal()).toBe(15);
  });

  it('deve formatar o rotulo da semana e mes corretamente', () => {
    expect(component.formatarRotuloSemana('2026-W36')).toContain('Semana 36');
    expect(component.formatarRotuloMes('2026-09')).toContain('setembro');
  });

  it('deve recuperar nome e cor da categoria', () => {
    expect(component.nomeDaCategoria(1)).toBe('Faculdade');
    expect(component.corDaCategoria(1)).toBe('#3b82f6');
    expect(component.nomeDaCategoria(99)).toBe('Categoria #99');
    expect(component.corDaCategoria(99)).toBe('#6c757d');
  });

  it('deve permitir trocar de periodo para semana e ano', () => {
    component.selecionarPeriodo('semana');
    expect(component.tipoPeriodo()).toBe('semana');

    component.selecionarPeriodo('ano');
    expect(component.tipoPeriodo()).toBe('ano');
  });

  it('deve exibir erro caso o servico falhe', () => {
    vi.spyOn(relatorioService, 'gerar').mockReturnValue(throwError(() => new Error('Falha')));
    component.carregarRelatorio();
    expect(component.erro()).toBe('Não foi possível carregar o relatório para este período.');
  });
});
