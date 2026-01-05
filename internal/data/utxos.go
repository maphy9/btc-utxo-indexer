package data

import "context"

type UtxosQ interface {
	GetActiveByAddress(ctx context.Context, address string) ([]Utxo, error)
	Spend(txHash string, txPos int, spentTxHash string) error
	Insert(utxo Utxo) (*Utxo, error)
}

type Utxo struct {
	Address     string  `db:"address" structs:"address" json:"-"`
	TxHash      string  `db:"tx_hash" structs:"tx_hash" json:"tx_hash"`
	TxPos       int     `db:"tx_pos" structs:"tx_pos" json:"tx_pos"`
	SpentTxHash *string `db:"spent_tx_hash" structs:"spent_tx_hash" json:"spent_tx_hash"`
	Value       int64   `db:"value" structs:"value" json:"value"`
}
