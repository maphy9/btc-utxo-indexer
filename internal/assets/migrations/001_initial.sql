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
  address text UNIQUE NOT NULL,
  user_id bigint REFERENCES users (id) ON DELETE CASCADE NOT NULL,
  created_at timestamptz DEFAULT current_timestamp,
  UNIQUE(user_id, address)
);

CREATE TABLE IF NOT EXISTS utxos
(
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  address text REFERENCES tracked_addresses (address) ON DELETE CASCADE NOT NULL,
  txid text NOT NULL,
  vout integer NOT NULL,
  amount bigint NOT NULL,
  block_height integer NOT NULL,
  block_hash text NOT NULL,
  UNIQUE(txid, vout)
);

-- +migrate Down
DROP TABLE IF EXISTS utxos;

DROP TABLE IF EXISTS tracked_addresses CASCADE;

DROP TABLE IF EXISTS users CASCADE;