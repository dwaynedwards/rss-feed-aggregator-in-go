-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id bigint GENERATED ALWAYS AS IDENTITY,
  email text NOT NULL,
  password text NOT NULL,
  name text NOT NULL,
  created_at timestamp NOT NULL DEFAULT now(),
  modefied_at timestamp NOT NULL DEFAULT now(),
  CONSTRAINT user_pk PRIMARY KEY (id),
  CONSTRAINT unique_email_address UNIQUE (email),
  CONSTRAINT check_email_length CHECK (char_length(email)<=256),
  CONSTRAINT check_name_length CHECK (char_length(name)<=50)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_modified_at_column() RETURNS TRIGGER AS $$
  BEGIN
    NEW.modefied_at = clock_timestamp();
    RETURN NEW;
  END;
$$ LANGUAGE 'plpgsql';
-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE TRIGGER update_users_modified_at_timestamp
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE PROCEDURE update_modified_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_users_modified_at_timestamp ON users;
-- +goose StatementEnd
-- +goose StatementBegin
DROP FUNCTION IF EXISTS update_modified_at_column;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
