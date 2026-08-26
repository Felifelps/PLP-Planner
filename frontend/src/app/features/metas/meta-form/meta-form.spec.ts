import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { beforeEach, describe, expect, it } from 'vitest';
import { MetaForm } from './meta-form';

describe('MetaForm', () => {
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MetaForm],
      providers: [provideHttpClient(), provideHttpClientTesting(), provideRouter([])],
    }).compileComponents();
  });

  it('should create', () => {
    const fixture = TestBed.createComponent(MetaForm);
    expect(fixture.componentInstance).toBeTruthy();
  });
});
