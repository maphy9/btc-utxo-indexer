package data

import "context"

type UtxosQ interface {
	GetByAddress(ctx context.Context, address string) ([]Utxo, error)

	InsertMany(ctx context.Context, utxos []Utxo) ([]Utxo, error)
}

type Utxo struct {
	ID          int64  `db:"id" structs:"-" json:"-"`
	AddressID   int64  `db:"address_id" structs:"address_id" json:"-"`
	TxID        string `db:"txid" structs:"txid" json:"txid"`
	Vout        uint   `db:"vout" structs:"vout" json:"vout"`
	Value       int64  `db:"value" structs:"value" json:"value"`
	BlockHeight int    `db:"block_height" structs:"block_height" json:"block_height"`
	BlockHash   string `db:"block_hash" structs:"block_hash" json:"block_hash"`
}
