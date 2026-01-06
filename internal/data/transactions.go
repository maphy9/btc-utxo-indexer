package data

import "context"

type TransactionsQ interface {
	InsertBatch(ctx context.Context, txs []Transaction) error
}

type Transaction struct {
	TxHash string `db:"tx_hash" structs:"tx_hash" json:"tx_hash"`
	Height int    `db:"height" structs:"height" json:"height"`
}
