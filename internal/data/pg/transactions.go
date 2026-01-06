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

func newTransactionsQ(db *pgdb.DB) data.TransactionsQ {
	return &transactionsQ{
		db:  db,
		sql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

type transactionsQ struct {
	db  *pgdb.DB
	sql squirrel.StatementBuilderType
}

func (m *transactionsQ) InsertBatch(ctx context.Context, txs []data.Transaction) error {
	if len(txs) == 0 {
		return nil
	}

	query := m.sql.Insert(transactionsTableName).
		Columns("tx_hash", "height").
		Suffix("ON CONFLICT DO NOTHING")

	for _, tx := range txs {
		query = query.Values(tx.TxHash, tx.Height)
	}

	return m.db.ExecContext(ctx, query)
}
