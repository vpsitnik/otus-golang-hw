-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    duration INTERVAL NOT NULL,
    description VARCHAR(300) NOT NULL,
    owner VARCHAR(50) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
