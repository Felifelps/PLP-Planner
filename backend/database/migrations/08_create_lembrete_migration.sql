-- +goose Up

CREATE TABLE IF NOT EXISTS lembretes (
    id BIGSERIAL PRIMARY KEY,

    descricao TEXT NOT NULL,

    tipo VARCHAR(50) NOT NULL
        CHECK (
            tipo IN (
                'reunião',
                'ligação',
                'compra',
                'estudo',
                'exercício',
                'entrega'
            )
        ),

    data DATE NOT NULL,

    horario TIME NOT NULL,

    recorrente BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down

DROP TABLE IF EXISTS lembretes;