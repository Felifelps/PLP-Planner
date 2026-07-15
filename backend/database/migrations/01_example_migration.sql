-- +goose Up

CREATE TABLE IF NOT EXISTS exemplos (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL
);

INSERT INTO exemplos (nome)
VALUES ('Exemplo inserido na migração de exemplo');

-- +goose Down

DROP TABLE IF EXISTS exemplos;