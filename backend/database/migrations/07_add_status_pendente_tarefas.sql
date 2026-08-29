-- +goose Up
ALTER TABLE tarefas DROP CONSTRAINT IF EXISTS tarefas_status_check;

ALTER TABLE tarefas ADD CONSTRAINT tarefas_status_check 
    CHECK (
        status IN (
            'pendente',
            'executada',
            'parcialmente executada',
            'cancelada',
            'adiada'
        )
    );

-- +goose Down
ALTER TABLE tarefas DROP CONSTRAINT IF EXISTS tarefas_status_check;

ALTER TABLE tarefas ADD CONSTRAINT tarefas_status_check 
    CHECK (
        status IN (
            'executada',
            'parcialmente executada',
            'cancelada',
            'adiada'
        )
    );