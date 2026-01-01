package data

import "context"

type TransactionsQ interface {
	GetByAddress(ctx context.Context, address string) ([]Transaction, error)

	InsertMany(ctx context.Context, transactions []Transaction) ([]Transaction, error)
}

type Transaction struct {
	TxHash string `db:"tx_hash" structs:"tx_hash" json:"tx_hash"`
	Height int    `db:"height" structs:"height" json:"height"`
}
