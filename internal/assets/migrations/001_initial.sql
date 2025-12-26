-- +migrate Up

CREATE TABLE IF NOT EXISTS users
(
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  username text UNIQUE NOT NULL,
  password_hash text NOT NULL,
  created_at timestamptz DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS addresses
(
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  address text NOT NULL,
  user_id bigint REFERENCES users (id) ON DELETE CASCADE NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS addresses CASCADE;

DROP TABLE IF EXISTS users CASCADE;