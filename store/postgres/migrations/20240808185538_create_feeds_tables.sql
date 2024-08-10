-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS feeds (
  id bigint GENERATED ALWAYS AS IDENTITY,
  url text NOT NULL,
  enabled boolean NOT NULL DEFAULT TRUE,
  deleted boolean NOT NULL DEFAULT FALSE,
  created_at timestamp NOT NULL,
  modified_at timestamp NOT NULL,
  last_synced_at timestamp NOT NULL,
  CONSTRAINT pk_feeds PRIMARY KEY (id),
  CONSTRAINT unique_url_address UNIQUE (url)
);

CREATE TABLE IF NOT EXISTS user_feeds (
  user_id bigint NOT NULL,
  feed_id bigint NOT NULL,
  name text NOT NULL,
  enabled boolean NOT NULL DEFAULT TRUE,
  deleted boolean NOT NULL DEFAULT FALSE,
  created_at timestamp NOT NULL,
  modified_at timestamp NOT NULL,
  CONSTRAINT pk_user_feed PRIMARY KEY (user_id, feed_id),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT fk_feed FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE,
  CONSTRAINT check_name_length CHECK (char_length(name)<=50)
);

CREATE TABLE IF NOT EXISTS feed_channels (
  id bigint GENERATED ALWAYS AS IDENTITY,
  feed_id bigint NOT NULL,
  title text NOT NULL,
  desciption text NOT NULL,
  link text NOT NULL,
  created_at timestamp NOT NULL,
  modified_at timestamp NOT NULL,
  CONSTRAINT pk_feed_channels PRIMARY KEY (id),
  CONSTRAINT fk_feed FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS feed_channel_items (
  id bigint GENERATED ALWAYS AS IDENTITY,
  feed_channel_id bigint NOT NULL,
  title text NOT NULL,
  desciption text NOT NULL,
  link text NOT NULL,
  created_at timestamp NOT NULL,
  modified_at timestamp NOT NULL,
  CONSTRAINT pk_feed_channel_items PRIMARY KEY (id),
  CONSTRAINT fk_feed_channel FOREIGN KEY (feed_channel_id) REFERENCES feed_channels (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS feed_channel_items;

DROP TABLE IF EXISTS feed_channels;

DROP TABLE IF EXISTS user_feeds;

DROP TABLE IF EXISTS feeds;
-- +goose StatementEnd
