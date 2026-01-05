package data

import "context"

type TransactionsQ interface {
	Exists(ctx context.Context, txHash string) (bool, error)
	Insert(ctx context.Context, tx Transaction) (*Transaction, error)
}

type Transaction struct {
	TxHash string `db:"tx_hash" structs:"tx_hash" json:"tx_hash"`
	Height int    `db:"height" structs:"height" json:"height"`
}
