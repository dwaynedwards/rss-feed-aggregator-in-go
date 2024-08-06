-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id bigint GENERATED ALWAYS AS IDENTITY,
  name text NOT NULL,
  created_at timestamp NOT NULL,
  modified_at timestamp NOT NULL,
  CONSTRAINT pk_users PRIMARY KEY (id),
  CONSTRAINT check_name_length CHECK (char_length(name)<=50)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
