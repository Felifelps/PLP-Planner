-- +goose Up

CREATE TABLE IF NOT EXISTS categorias (
    id BIGSERIAL PRIMARY KEY,

    nome VARCHAR(100) NOT NULL UNIQUE,

    cor VARCHAR(7) NOT NULL
        CHECK (cor ~ '^#[0-9A-Fa-f]{6}$')
);

INSERT INTO categorias (nome, cor) VALUES
    ('Trabalho', '#4C6EF5'),
    ('Estudos', '#7048E8'),
    ('Saúde', '#12B886'),
    ('Lazer', '#F59F00'),
    ('Casa', '#F76707'),
    ('Pessoal', '#E64980');

-- +goose Down

DROP TABLE IF EXISTS categorias;
