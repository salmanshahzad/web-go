-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user"(
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "user";
-- +goose StatementEnd
