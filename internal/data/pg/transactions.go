package pg

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
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

func (m *transactionsQ) Exists(ctx context.Context, txHash string) (bool, error) {
	query := m.sql.Select("COUNT(*)").
		From(transactionsTableName).
		Where("tx_hash = ?", txHash)

	var result int
	err := m.db.GetContext(ctx, &result, query)
	return result > 0, err
}

func (m *transactionsQ) Insert(ctx context.Context, tx data.Transaction) (*data.Transaction, error) {
	clauses := structs.Map(tx)
	query := m.sql.Insert(transactionsTableName).
		SetMap(clauses).
		Suffix("RETURNING *")

	var result data.Transaction
	err := m.db.GetContext(ctx, &result, query)
	return &result, err
}
