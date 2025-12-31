package pg

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	transactionsTableName = "transactions"
)

func NewTransactionsQ(db *pgdb.DB) data.TransactionsQ {
	return &transactionsQ{
		db: db,
		sql: squirrel.StatementBuilder,
	}
}

type transactionsQ struct {
	db *pgdb.DB
	sql squirrel.StatementBuilderType
}

func (m *transactionsQ) GetByAddress(ctx context.Context, address string) ([]data.Transaction, error) {
	query := m.sql.Select("t.*").
		Distinct().
		From(transactionsTableName + " t").
		Join(utxosTableName + " u ON t.txid = u.txid").
		Where("u.address = ?", address).
		PlaceholderFormat(squirrel.Dollar)

	var result []data.Transaction
	err := m.db.SelectContext(ctx, &result, query)
	return result, err
}

func (m *transactionsQ) InsertMany(ctx context.Context, transactions []data.Transaction) ([]data.Transaction, error) {
	if len(transactions) == 0 {
		return nil, nil
	}

	query := m.sql.Insert(transactionsTableName).
		Columns("txid", "block_height")

	for _, transation := range transactions {
		query = query.Values(transation.TxID, transation.BlockHeight)
	}

	query = query.Suffix("ON CONFLICT (txid) DO NOTHING RETURNING *")

	var result []data.Transaction
	err := m.db.SelectContext(ctx, &result, query)
	return result, err
}
