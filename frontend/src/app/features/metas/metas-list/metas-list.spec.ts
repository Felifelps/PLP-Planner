import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { beforeEach, describe, expect, it } from 'vitest';
import { MetasList } from './metas-list';

describe('MetasList', () => {
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MetasList],
      providers: [provideHttpClient(), provideHttpClientTesting(), provideRouter([])],
    }).compileComponents();
  });

  it('should create', () => {
    const fixture = TestBed.createComponent(MetasList);
    expect(fixture.componentInstance).toBeTruthy();
  });
});
