-- +migrate Up

CREATE TABLE IF NOT EXISTS users
(
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  username text UNIQUE NOT NULL,
  password_hash text NOT NULL,
  refresh_token text NOT NULL DEFAULT '',
  created_at timestamptz DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS tracked_addresses
(
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  addr text NOT NULL,
  user_id bigint REFERENCES users (id) ON DELETE CASCADE NOT NULL,
  created_at timestamptz DEFAULT current_timestamp,
  UNIQUE(user_id, addr)
);

-- +migrate Down
DROP TABLE IF EXISTS tracked_addresses CASCADE;

DROP TABLE IF EXISTS users CASCADE;