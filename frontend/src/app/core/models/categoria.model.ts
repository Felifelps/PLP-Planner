export interface Categoria {
  id: number;
  nome: string;
  cor: string;
}

export type CategoriaPayload = Omit<Categoria, 'id'>;
