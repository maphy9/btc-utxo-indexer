package data

import "context"

type UtxosQ interface {
	GetByAddress(ctx context.Context, address string) ([]Utxo, error)
	Spend(txHash string, txPos int, spentHeight int) error
	Insert(utxo Utxo) (*Utxo, error)
}

type Utxo struct {
	Address       string `db:"address" structs:"address" json:"-"`
	TxHash        string `db:"tx_hash" structs:"tx_hash" json:"tx_hash"`
	TxPos         int    `db:"tx_pos" structs:"tx_pos" json:"tx_pos"`
	Value         int64  `db:"value" structs:"value" json:"value"`
	CreatedHeight int    `db:"created_height" structs:"created_height" json:"created_height"`
	SpentHeight   *int   `db:"spent_height" structs:"spent_height" json:"spent_height"`
}
