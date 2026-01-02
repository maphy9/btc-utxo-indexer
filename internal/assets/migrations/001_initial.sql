-- +migrate Up

CREATE TABLE IF NOT EXISTS users
(
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  username text UNIQUE NOT NULL,
  password_hash text NOT NULL,
  refresh_token text NOT NULL DEFAULT '',
  created_at timestamptz DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS addresses
(
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  address text UNIQUE NOT NULL,
  status text NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS user_addresses
(
  address_id bigint REFERENCES addresses (id) NOT NULL,
  user_id bigint REFERENCES users (id) ON DELETE CASCADE NOT NULL,
  UNIQUE (address_id, user_id)
);

CREATE TABLE IF NOT EXISTS headers (
  height integer PRIMARY KEY,
  hash text UNIQUE NOT NULL,
  parent_hash text NOT NULL,
  root text NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
  tx_hash text PRIMARY KEY,
  height integer REFERENCES headers (height) ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS utxos
(
  address text REFERENCES addresses (address) ON DELETE CASCADE NOT NULL,
  tx_hash text REFERENCES transactions (tx_hash) ON DELETE CASCADE NOT NULL,
  tx_pos integer NOT NULL,
  value bigint NOT NULL,
  created_height integer REFERENCES headers (height) ON DELETE CASCADE NOT NULL,
  spent_height integer REFERENCES headers (height) ON DELETE SET NULL,
  PRIMARY KEY (tx_hash, tx_pos)
);

CREATE INDEX idx_utxos_address ON utxos (address);

-- +migrate Down
DROP TABLE IF EXISTS utxos CASCADE;

DROP TABLE IF EXISTS transactions CASCADE;

DROP TABLE IF EXISTS headers CASCADE;

DROP TABLE IF EXISTS user_addresses CASCADE;

DROP TABLE IF EXISTS addresses CASCADE;

DROP TABLE IF EXISTS users CASCADE;