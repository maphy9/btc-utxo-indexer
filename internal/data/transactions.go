package data

import "context"

type TransactionsQ interface {
	GetByAddress(ctx context.Context, address string) ([]Transaction, error)

	InsertMany(ctx context.Context, transactions []Transaction) ([]Transaction, error)
}

type Transaction struct {
	TxID        string `db:"txid" structs:"txid" json:"txid"`
	BlockHeight int    `db:"block_height" structs:"block_height" json:"block_height"`
}
