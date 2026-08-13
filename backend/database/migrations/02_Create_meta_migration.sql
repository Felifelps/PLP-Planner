-- +goose Up

CREATE TABLE IF NOT EXISTS metas (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    descricao TEXT,
    categoria_id BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL,
    data_inicio DATE NOT NULL,
    data_fim DATE NOT NULL
);

-- +goose Down

DROP TABLE IF EXISTS metas;