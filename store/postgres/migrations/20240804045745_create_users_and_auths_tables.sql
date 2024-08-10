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

CREATE TABLE IF NOT EXISTS auths (
  id bigint GENERATED ALWAYS AS IDENTITY,
  user_id bigint NOT NULL,
  email text NOT NULL,
  password text NOT NULL,
  enabled boolean NOT NULL DEFAULT TRUE,
  deleted boolean NOT NULL DEFAULT FALSE,
  created_at timestamp NOT NULL,
  modified_at timestamp NOT NULL,
  last_signed_in_at timestamp NOT NULL,
  CONSTRAINT pk_auths PRIMARY KEY (id),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT unique_email_address UNIQUE (email),
  CONSTRAINT check_email_length CHECK (char_length(email)<=256)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS auths;

DROP TABLE IF EXISTS users;
-- +goose StatementEnd
