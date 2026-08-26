-- +goose Up

ALTER TABLE metas
    ADD CONSTRAINT metas_categoria_id_fkey
        FOREIGN KEY (categoria_id) REFERENCES categorias (id);

-- +goose Down

ALTER TABLE metas
    DROP CONSTRAINT metas_categoria_id_fkey;
