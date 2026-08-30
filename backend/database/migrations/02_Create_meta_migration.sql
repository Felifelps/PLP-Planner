-- +goose Up

CREATE TABLE IF NOT EXISTS metas (
    id BIGSERIAL PRIMARY KEY,

    nome VARCHAR(255) NOT NULL,

    descricao TEXT,

    categoria_id BIGINT NOT NULL,

    status VARCHAR(50) NOT NULL
        CHECK (
            status IN (
                'cumprida',
                'parcialmente cumprida',
                'não cumprida'
            )
        ),

    periodo VARCHAR(20) NOT NULL
        CHECK (
            periodo IN (
                'diario',
                'semanal',
                'mensal',
                'anual'
            )
        ),

    data_inicio DATE NOT NULL,

    data_fim DATE NOT NULL,

    CONSTRAINT metas_periodo_datas_check
        CHECK (data_inicio <= data_fim)
);

-- +goose Down

DROP TABLE IF EXISTS metas;