-- +goose Up

INSERT INTO metas (nome, descricao, categoria_id, status, periodo, data_inicio, data_fim) VALUES
    ('Estudar para a prova de PLP', 'Revisar os capítulos 3 e 4 antes da prova', 2, 'não cumprida', 'semanal', '2026-08-24', '2026-08-30'),
    ('Correr 3x na semana', 'Manter a rotina de corridas curtas', 3, 'parcialmente cumprida', 'semanal', '2026-08-24', '2026-08-30'),
    ('Organizar o home office', 'Arrumar mesa e cabos', 5, 'cumprida', 'semanal', '2026-08-17', '2026-08-23'),
    ('Ler um livro técnico', 'Terminar a leitura do livro do mês', 2, 'não cumprida', 'mensal', '2026-08-01', '2026-08-31'),
    ('Fechar as tarefas do projeto PLP-Planner', 'Concluir o módulo de tarefas e categorias', 1, 'parcialmente cumprida', 'mensal', '2026-08-01', '2026-08-31'),
    ('Economizar 10% do salário', 'Guardar parte do salário do mês', 1, 'não cumprida', 'mensal', '2026-07-01', '2026-07-31'),
    ('Viajar nas férias de fim de ano', 'Planejar destino e orçamento', 4, 'não cumprida', 'anual', '2026-01-01', '2026-12-31'),
    ('Manter rotina de exercícios', 'Treinar ao menos 3x por semana durante o ano', 3, 'parcialmente cumprida', 'anual', '2026-01-01', '2026-12-31'),
    ('Passar mais tempo com a família', 'Reservar finais de semana para a família', 6, 'cumprida', 'anual', '2026-01-01', '2026-12-31'),
    ('Planejar a próxima semana de estudos', 'Organizar o cronograma da semana seguinte', 2, 'não cumprida', 'semanal', '2026-08-31', '2026-09-06');

INSERT INTO tarefas (descricao, categoria_id, data, horario_inicio, duracao, turno, status, prioridade) VALUES
    ('Preparar apresentação do projeto', 1, '2026-08-25', '14:00', '1h', NULL, 'executada', 'alta'),
    ('Revisar slides de PLP', 2, '2026-08-26', '19:00', '1h', NULL, 'adiada', 'alta'),
    ('Corrida matinal', 3, '2026-08-26', NULL, NULL, 'manhã', 'executada', 'média'),
    ('Reunião de equipe', 1, '2026-08-26', '09:00', '30min', NULL, 'executada', 'alta'),
    ('Lavar roupas', 5, '2026-08-26', NULL, NULL, 'tarde', 'parcialmente executada', 'baixa'),
    ('Assistir um filme', 4, '2026-08-27', NULL, NULL, 'noite', 'cancelada', 'baixa'),
    ('Ligar para a família', 6, '2026-08-27', '20:00', '30min', NULL, 'adiada', 'média'),
    ('Estudar Angular', 2, '2026-08-28', NULL, NULL, 'noite', 'adiada', 'alta'),
    ('Academia', 3, '2026-08-29', NULL, NULL, 'tarde', 'adiada', 'média'),
    ('Organizar a despensa', 5, '2026-08-30', NULL, NULL, 'manhã', 'adiada', 'baixa');

-- +goose Down

DELETE FROM tarefas WHERE descricao IN (
    'Preparar apresentação do projeto',
    'Revisar slides de PLP',
    'Corrida matinal',
    'Reunião de equipe',
    'Lavar roupas',
    'Assistir um filme',
    'Ligar para a família',
    'Estudar Angular',
    'Academia',
    'Organizar a despensa'
);

DELETE FROM metas WHERE nome IN (
    'Estudar para a prova de PLP',
    'Correr 3x na semana',
    'Organizar o home office',
    'Ler um livro técnico',
    'Fechar as tarefas do projeto PLP-Planner',
    'Economizar 10% do salário',
    'Viajar nas férias de fim de ano',
    'Manter rotina de exercícios',
    'Passar mais tempo com a família',
    'Planejar a próxima semana de estudos'
);
