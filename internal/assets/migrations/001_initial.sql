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
  address_id bigint REFERENCES addresses (id) ON DELETE CASCADE NOT NULL,
  user_id bigint REFERENCES users (id) ON DELETE CASCADE NOT NULL,
  UNIQUE (address_id, user_id)
);

CREATE TABLE IF NOT EXISTS headers (
  height integer PRIMARY KEY,
  hash text UNIQUE NOT NULL,
  parent_hash text NOT NULL,
  root text NOT NULL,
  created_at timestamp
);

CREATE TABLE IF NOT EXISTS transactions (
  tx_hash text PRIMARY KEY,
  height integer REFERENCES headers (height) ON DELETE CASCADE NOT NULL
);

CREATE TABLE transaction_outputs (
  tx_hash text REFERENCES transactions (tx_hash) ON DELETE CASCADE NOT NULL,
  output_index integer NOT NULL,
  value bigint NOT NULL,
  address text,
  script_hex text NOT NULL,
  spent_by_tx_hash text REFERENCES transactions (tx_hash) ON DELETE SET NULL,
  PRIMARY KEY (tx_hash, output_index)
);

CREATE TABLE transaction_inputs (
  tx_hash text REFERENCES transactions (tx_hash) ON DELETE CASCADE NOT NULL,
  input_index integer NOT NULL,
  prev_tx_hash text NOT NULL,
  prev_output_index integer NOT NULL,
  PRIMARY KEY (tx_hash, input_index)
);

CREATE INDEX idx_tx_outputs_address ON transaction_outputs (address);

CREATE INDEX idx_tx_outputs_unspent ON transaction_outputs (address) 
WHERE spent_by_tx_hash IS NULL;

CREATE INDEX idx_tx_inputs_prev ON transaction_inputs (prev_tx_hash, prev_output_index);

CREATE INDEX idx_tx_outputs_spent_by ON transaction_outputs (spent_by_tx_hash);

-- +migrate Down
DROP TABLE IF EXISTS transaction_inputs CASCADE;

DROP TABLE IF EXISTS transaction_outputs CASCADE;

DROP TABLE IF EXISTS transactions CASCADE;

DROP TABLE IF EXISTS headers CASCADE;

DROP TABLE IF EXISTS user_addresses CASCADE;

DROP TABLE IF EXISTS addresses CASCADE;

DROP TABLE IF EXISTS users CASCADE;