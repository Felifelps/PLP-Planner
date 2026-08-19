package models

import "testing"

func TestCategoriaValidate(t *testing.T) {
	testes := []struct {
		nome      string
		categoria Categoria
		erro      bool
	}{
		{
			nome: "categoria válida",
			categoria: Categoria{
				Nome: "Trabalho",
				Cor:  "#4C6EF5",
			},
			erro: false,
		},
		{
			nome: "cor sem #",
			categoria: Categoria{
				Nome: "Trabalho",
				Cor:  "4C6EF5",
			},
			erro: true,
		},
		{
			nome: "cor com tamanho inválido",
			categoria: Categoria{
				Nome: "Trabalho",
				Cor:  "#FFF",
			},
			erro: true,
		},
		{
			nome: "cor vazia",
			categoria: Categoria{
				Nome: "Trabalho",
				Cor:  "",
			},
			erro: true,
		},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			err := tt.categoria.Validate()

			if (err != nil) != tt.erro {
				t.Fatalf(
					"Validate() erro = %v; esperado erro = %v",
					err,
					tt.erro,
				)
			}
		})
	}
}

func TestCorValida(t *testing.T) {
	cores := []string{
		"#000000",
		"#FFFFFF",
		"#4c6ef5",
		"#F59F00",
	}

	for _, cor := range cores {
		t.Run(cor, func(t *testing.T) {
			if !CorValida(cor) {
				t.Fatalf("cor %q deveria ser válida", cor)
			}
		})
	}
}
