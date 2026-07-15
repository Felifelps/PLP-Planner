# PLP-Planner

O PLP-Planner é um planner desenvolvido para a disciplina de PLP.

A aplicação é composta por um frontend em Angular, um backend em Go e um banco de dados PostgreSQL.

Todo o ambiente é executado com Docker. Não é necessário instalar Go, Node.js, Angular CLI ou PostgreSQL diretamente na máquina.

## Documentação do backend

As orientações para adicionar models, repositories, services, handlers, rotas e migrations estão disponíveis no README do backend:

[Acessar o padrão de desenvolvimento do backend](backend/README.md)

## Requisitos

- Docker;
- Docker Compose.

## Executando o projeto

Na raiz do projeto, execute:

```bash
docker compose up --build
```

O primeiro início pode demorar alguns minutos, pois o Docker precisará construir as imagens e instalar as dependências.

Para executar os containers em segundo plano:

```bash
docker compose up --build -d
```

## Acessos

Depois que os containers estiverem prontos, os serviços estarão disponíveis nos seguintes endereços:

| Serviço | Endereço |
| --- | --- |
| Frontend | http://localhost:4200 |
| Backend | http://localhost:8080 |
| PostgreSQL | `localhost:5432` |

O endereço `localhost` deve ser utilizado no navegador e em aplicações executadas na máquina host. Dentro da rede do Docker, os containers se comunicam pelos nomes dos serviços, como `backend` e `postgres`.

## Desenvolvimento

As pastas locais do frontend e do backend são montadas dentro dos containers. Assim, alterações no código são refletidas durante o desenvolvimento sem a necessidade de reconstruir as imagens a cada mudança.

- o backend utiliza live reload;
- o frontend utiliza o servidor de desenvolvimento do Angular;
- os dados do PostgreSQL são preservados no volume `postgres_data`.

## Comandos úteis

Exibir os logs:

```bash
docker compose logs -f
```

Exibir os logs de apenas um serviço:

```bash
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f postgres
```

Parar os containers:

```bash
docker compose down
```

Parar os containers e apagar os dados do banco:

```bash
docker compose down -v
```

> O último comando remove o volume do PostgreSQL e apaga os dados armazenados. Use-o somente quando quiser reiniciar o banco.