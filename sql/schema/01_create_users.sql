-- +goose Up
CREATE TABLE "user"(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT(
        encode(sha256(random()::text::bytea), 'hex'))
);
-- +goose Down
DROP TABLE "user";