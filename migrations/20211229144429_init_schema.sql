-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA if not exists shortener;
CREATE TABLE if not exists shortener.url (
	id BIGSERIAL primary key,
	full_url text,
	user_token text,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA shortener CASCADE;
-- +goose StatementEnd
