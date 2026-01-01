package data

import "context"

type UtxosQ interface {
	GetByAddress(ctx context.Context, address string) ([]Utxo, error)

	InsertMany(ctx context.Context, utxos []Utxo) ([]Utxo, error)
}

type Utxo struct {
	Address string `db:"address" structs:"address" json:"-"`
	TxHash  string `db:"tx_hash" structs:"tx_hash" json:"tx_hash"`
	TxPos   uint   `db:"tx_pos" structs:"tx_pos" json:"tx_pos"`
	Value   int64  `db:"value" structs:"value" json:"value"`
	Height  int    `db:"height" structs:"height" json:"height"`
}
