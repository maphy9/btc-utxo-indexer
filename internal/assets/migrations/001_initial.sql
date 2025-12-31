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
  address text UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS user_addresses
(
  address_id bigint REFERENCES addresses (id) NOT NULL,
  user_id bigint REFERENCES users (id) ON DELETE CASCADE NOT NULL,
  UNIQUE (address_id, user_id)
);

-- CREATE TABLE IF NOT EXISTS blocks (
--   height integer PRIMARY KEY,
--   hash text UNIQUE NOT NULL,
--   parent_hash text NOT NULL,
--   timestamp timestamptz NOT NULL 
-- );

-- CREATE TABLE IF NOT EXISTS transactions (
--   id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--   txid text UNIQUE NOT NULL,
--   block_height integer REFERENCES blocks (height) ON DELETE CASCADE NOT NULL
-- );

CREATE TABLE IF NOT EXISTS utxos
(
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  address text REFERENCES addresses (address) ON DELETE CASCADE NOT NULL,
  txid text NOT NULL,  -- REFERENCES transactions (txid) ON DELETE CASCADE
  block_height integer NOT NULL,  -- REFERENCES blocks (height) ON DELETE CASCADE
  block_hash text NOT NULL, -- REFERENCES blocks (hash) ON DELETE CASCADE
  vout integer NOT NULL,
  value bigint NOT NULL,
  UNIQUE (txid, vout)
);

-- +migrate Down
DROP TABLE IF EXISTS utxos CASCADE;

DROP TABLE IF EXISTS transactions CASCADE;

DROP TABLE IF EXISTS blocks CASCADE;

DROP TABLE IF EXISTS user_addresses CASCADE;

DROP TABLE IF EXISTS addresses CASCADE;

DROP TABLE IF EXISTS users CASCADE;