-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stored_items
(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    type varchar(30),
    description text,
    value text,
    userId BIGINT references users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stored_items;
-- +goose StatementEnd
