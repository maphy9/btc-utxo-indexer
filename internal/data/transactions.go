package data

import (
	"context"
	"time"
)

type TransactionsQ interface {
	GetUtxosByAddress(ctx context.Context, address string) ([]TransactionOutput, error)
	GetAddressBalance(ctx context.Context, address string) (int64, error)
	GetAddressTransactions(ctx context.Context, address string) ([]AddressTransaction, error)
	InsertTransactionsBatch(ctx context.Context, txs []Transaction) error
	InsertTransactionOutputsBatch(ctx context.Context, txs []TransactionOutput) error
	InsertTransactionInputsBatch(ctx context.Context, txs []TransactionInput) error
	SpendTransactionOutputs(ctx context.Context, txs []TransactionInput) error
}

type Transaction struct {
	TxHash    string    `db:"tx_hash"`
	CreatedAt time.Time `db:"created_at"`
	Height    int       `db:"height"`
}

type TransactionOutput struct {
	TxHash        string  `db:"tx_hash"`
	Index         int     `db:"output_index"`
	Value         int64   `db:"value"`
	Address       string  `db:"address"`
	ScriptHex     string  `db:"script_hex"`
	SpentByTxHash *string `db:"spent_by_tx_hash"`
}

type TransactionInput struct {
	TxHash     string `db:"tx_hash"`
	Index      int    `db:"input_index"`
	PrevTxHash string `db:"prev_tx_hash"`
	PrevIndex  int    `db:"prev_output_index"`
}

type AddressTransaction struct {
	TxHash    string    `db:"tx_hash" json:"tx_hash"`
	Height    string    `db:"height" json:"height"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Delta     int64     `db:"delta" json:"delta"`
}
