-- +goose Up

CREATE TABLE IF NOT EXISTS tarefas (
    id BIGSERIAL PRIMARY KEY,

    descricao TEXT NOT NULL,

    categoria_id BIGINT NOT NULL
        REFERENCES categorias (id),

    data DATE NOT NULL,

    horario_inicio VARCHAR(5),

    duracao VARCHAR(10)
        CHECK (
            duracao IS NULL
            OR duracao IN ('30min', '1h')
        ),

    turno VARCHAR(10)
        CHECK (
            turno IS NULL
            OR turno IN ('manhã', 'tarde', 'noite')
        ),

    status VARCHAR(30) NOT NULL
        CHECK (
            status IN (
                'executada',
                'parcialmente executada',
                'cancelada',
                'adiada'
            )
        ),

    prioridade VARCHAR(10) NOT NULL
        CHECK (
            prioridade IN ('baixa', 'média', 'alta')
        ),

    CONSTRAINT tarefas_horario_ou_turno_check CHECK (
        (horario_inicio IS NOT NULL AND duracao IS NOT NULL AND turno IS NULL)
        OR
        (horario_inicio IS NULL AND duracao IS NULL AND turno IS NOT NULL)
    )
);

-- +goose Down

DROP TABLE IF EXISTS tarefas;
