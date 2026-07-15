# Tutorial — Adicionando uma nova funcionalidade ao backend

Este tutorial apresenta o padrão para adicionar uma nova funcionalidade ao backend Go do projeto.

O exemplo utiliza uma entidade genérica chamada `Recurso`. Substitua esse nome pela entidade real que será implementada.

## Fluxo das camadas

```text
Rota → Handler → Service → Repository → PostgreSQL
```

- `models`: define as estruturas de dados;
- `repositories`: realiza operações no banco;
- `services`: concentra as regras de negócio;
- `handlers`: trata requisições e respostas HTTP;
- `bootstrap`: inicializa as dependências e registra as rotas;
- `database/migrations`: versiona as alterações do banco.

## Estrutura

```text
backend/
├── bootstrap/
│   ├── dependencies.go
│   └── router.go
├── database/
│   ├── migrations/
│   ├── migrations.go
│   └── postgres.go
├── handlers/
├── models/
├── repositories/
├── services/
├── go.mod
├── go.sum
└── main.go
```

# 1. Criar a migration

Toda alteração na estrutura do banco deve ser feita por uma nova migration. Não crie ou altere tabelas manualmente.

Na raiz do projeto, execute:

```bash
docker compose exec backend \
  goose -dir database/migrations create create_recursos sql
```

Será criado um arquivo semelhante a:

```text
backend/database/migrations/20260715200000_create_recursos.sql
```

Preencha o arquivo:

```sql
-- +goose Up

CREATE TABLE recursos (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL
);

-- +goose Down

DROP TABLE IF EXISTS recursos;
```

- `Up` aplica a alteração;
- `Down` desfaz a alteração.

As migrations pendentes são executadas automaticamente quando o backend inicia.

# 2. Criar o model

Crie:

```text
backend/models/recurso.go
```

```go
package models

type Recurso struct {
	ID   int64  `json:"id"`
	Nome string `json:"nome"`
}
```

O model representa os dados utilizados pela aplicação. Ele não acessa o banco e não contém regras HTTP.

# 3. Criar o repository

Crie:

```text
backend/repositories/recurso.go
```

```go
package repositories

import (
	"context"

	"plp-planner/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RecursoRepository struct {
	db *pgxpool.Pool
}

func NewRecursoRepository(
	db *pgxpool.Pool,
) *RecursoRepository {
	return &RecursoRepository{db: db}
}

func (r *RecursoRepository) BuscarTodos(
	ctx context.Context,
) ([]models.Recurso, error) {
	rows, err := r.db.Query(
		ctx,
		"SELECT id, nome FROM recursos ORDER BY id",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recursos := make([]models.Recurso, 0)

	for rows.Next() {
		var recurso models.Recurso

		if err := rows.Scan(
			&recurso.ID,
			&recurso.Nome,
		); err != nil {
			return nil, err
		}

		recursos = append(recursos, recurso)
	}

	return recursos, rows.Err()
}
```

O repository deve conter apenas operações relacionadas à persistência dos dados.

# 4. Criar o service

Crie:

```text
backend/services/recurso.go
```

```go
package services

import (
	"context"

	"plp-planner/models"
)

type RecursoRepository interface {
	BuscarTodos(ctx context.Context) ([]models.Recurso, error)
}

type RecursoService struct {
	repository RecursoRepository
}

func NewRecursoService(
	repository RecursoRepository,
) *RecursoService {
	return &RecursoService{
		repository: repository,
	}
}

func (s *RecursoService) BuscarTodos(
	ctx context.Context,
) ([]models.Recurso, error) {
	return s.repository.BuscarTodos(ctx)
}
```

O service contém as regras de negócio. Ele depende de uma interface, facilitando testes e evitando dependência direta da implementação do repository.

# 5. Criar o handler

Crie:

```text
backend/handlers/recurso.go
```

```go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"plp-planner/services"
)

type RecursoHandler struct {
	service *services.RecursoService
}

func NewRecursoHandler(
	service *services.RecursoService,
) *RecursoHandler {
	return &RecursoHandler{
		service: service,
	}
}

func (h *RecursoHandler) BuscarTodos(
	w http.ResponseWriter,
	r *http.Request,
) {
	recursos, err := h.service.BuscarTodos(r.Context())
	if err != nil {
		log.Printf("Erro ao buscar recursos: %v", err)

		http.Error(
			w,
			"Erro ao buscar recursos",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(recursos); err != nil {
		log.Printf("Erro ao gerar resposta: %v", err)
	}
}
```

O handler deve:

1. ler os dados da requisição;
2. chamar o service;
3. converter erros para respostas HTTP;
4. devolver o resultado em JSON.

# 6. Registrar as dependências

Abra:

```text
backend/bootstrap/dependencies.go
```

Adicione o repository:

```go
type Repositories struct {
	Exemplo *repositories.ExemploRepository
	Recurso *repositories.RecursoRepository
}
```

Adicione o service:

```go
type Services struct {
	Exemplo *services.ExemploService
	Recurso *services.RecursoService
}
```

Adicione o handler:

```go
type Handlers struct {
	Exemplo *handlers.ExemploHandler
	Recurso *handlers.RecursoHandler
}
```

Inicialize o repository:

```go
func initializeRepositories(
	db *pgxpool.Pool,
) *Repositories {
	return &Repositories{
		Exemplo: repositories.NewExemploRepository(db),
		Recurso: repositories.NewRecursoRepository(db),
	}
}
```

Inicialize o service:

```go
func initializeServices(
	appRepositories *Repositories,
) *Services {
	return &Services{
		Exemplo: services.NewExemploService(
			appRepositories.Exemplo,
		),
		Recurso: services.NewRecursoService(
			appRepositories.Recurso,
		),
	}
}
```

Inicialize o handler:

```go
func initializeHandlers(
	appServices *Services,
) *Handlers {
	return &Handlers{
		Exemplo: handlers.NewExemploHandler(
			appServices.Exemplo,
		),
		Recurso: handlers.NewRecursoHandler(
			appServices.Recurso,
		),
	}
}
```

A ordem das dependências deve ser mantida:

```text
Banco → Repository → Service → Handler
```

# 7. Registrar as rotas

Abra:

```text
backend/bootstrap/router.go
```

Crie uma função para as rotas do novo recurso:

```go
func initializeRecursoRoutes(
	router *http.ServeMux,
	appHandlers *Handlers,
) {
	router.HandleFunc(
		"GET /recursos",
		appHandlers.Recurso.BuscarTodos,
	)
}
```

Chame essa função em `InitializeRouter`:

```go
func InitializeRouter(
	dependencies *Dependencies,
) *http.ServeMux {
	router := http.NewServeMux()

	initializeStatusRoutes(router)
	initializeExemploRoutes(router, dependencies.Handlers)
	initializeRecursoRoutes(router, dependencies.Handlers)

	return router
}
```

Mantenha as rotas de cada domínio em uma função própria para evitar que `InitializeRouter` fique extenso.

# 8. Reiniciar e validar

Formate e valide o código:

```bash
docker compose exec backend gofmt -w .
docker compose exec backend go test ./...
```

Reinicie o backend para aplicar a migration automaticamente:

```bash
docker compose restart backend
```

Se uma dependência ou migration nova não for reconhecida, reconstrua a imagem:

```bash
docker compose up -d --build backend
```

Confira os logs:

```bash
docker compose logs -f backend
```

Teste a rota:

```bash
curl http://localhost:8080/recursos
```

O retorno inicial esperado é:

```json
[]
```

# Adicionando alterações futuras ao banco

Nunca altere uma migration que já foi aplicada ou enviada ao repositório. Crie uma nova migration.

Exemplo para adicionar uma coluna:

```bash
docker compose exec backend \
  goose -dir database/migrations create add_descricao_to_recursos sql
```

```sql
-- +goose Up

ALTER TABLE recursos
ADD COLUMN descricao TEXT;

-- +goose Down

ALTER TABLE recursos
DROP COLUMN descricao;
```

# Arquivos Go no Git

Os arquivos abaixo devem ser versionados:

```text
go.mod
go.sum
```

O `go.mod` declara as dependências diretas do projeto. O `go.sum` registra os hashes utilizados para verificar as versões baixadas e garantir instalações reproduzíveis.

Portanto, não adicione `go.sum` ao `.gitignore`.

Um `.gitignore` básico para o backend pode conter:

```gitignore
# Binários
/bin/
*.exe
*.out

# Cobertura e testes
coverage.out

# Variáveis locais
.env

# Arquivos de editor
.idea/
.vscode/
```

# Checklist

Antes de concluir uma nova funcionalidade, confira:

- [ ] a migration possui `Up` e `Down`;
- [ ] o model foi criado;
- [ ] o repository foi criado e registrado no bootstrap;
- [ ] o service foi criado e registrado no bootstrap;
- [ ] o handler foi criado e registrado no bootstrap;
- [ ] as rotas foram registradas;
- [ ] `gofmt` foi executado;
- [ ] os testes passaram;
- [ ] a migration foi aplicada automaticamente;
- [ ] a rota foi testada;
- [ ] `go.mod` e `go.sum` estão versionados.
