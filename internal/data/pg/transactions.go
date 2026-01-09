package pg

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	transactionsTableName       = "transactions"
	transactionOutputsTableName = "transaction_outputs"
	transactionInputsTableName  = "transaction_inputs"
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

func (m *transactionsQ) GetUtxosByAddress(ctx context.Context, address string) ([]data.TransactionOutput, error) {
	query := m.sql.Select("*").
		From(transactionOutputsTableName).
		Where("address = ? AND spent_by_tx_hash IS NULL", address)

	var result []data.TransactionOutput
	err := m.db.SelectContext(ctx, &result, query)
	return result, err
}

func (m *transactionsQ) GetAddressBalance(ctx context.Context, address string) (int64, error) {
	query := m.sql.Select("COALESCE(SUM(value), 0)").
		From(transactionOutputsTableName).
		Where("address = ? AND spent_by_tx_hash IS NULL", address)

	var result int64
	err := m.db.GetContext(ctx, &result, query)
	return result, err
}

func (m *transactionsQ) GetAddressTransactions(ctx context.Context, address string) ([]data.AddressTransaction, error) {
	query := squirrel.Expr(
		fmt.Sprintf(
			`SELECT txs.tx_hash, txs.height, hdrs.created_at, SUM(delta) AS delta
			FROM (
				SELECT tx_hash, COALESCE(SUM(value), 0) AS delta
				FROM %s
				WHERE address = $1
				GROUP BY tx_hash

				UNION ALL
				
				SELECT spent_by_tx_hash, -COALESCE(SUM(value), 0)
				FROM %s
				WHERE address = $1
				GROUP BY spent_by_tx_hash
			) AS tmp
			JOIN %s AS txs ON tmp.tx_hash = txs.tx_hash
			JOIN %s AS hdrs ON hdrs.height = txs.height
			GROUP BY txs.tx_hash, txs.height, hdrs.created_at
			ORDER BY created_at`,
			transactionOutputsTableName,
			transactionOutputsTableName,
			transactionsTableName,
			headersTableName,
		),
		address,
	)

	var result []data.AddressTransaction
	err := m.db.SelectContext(ctx, &result, query)
	return result, err
}

func (m *transactionsQ) InsertTransactionsBatch(ctx context.Context, txs []data.Transaction) error {
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

func (m *transactionsQ) InsertTransactionOutputsBatch(ctx context.Context, outs []data.TransactionOutput) error {
	if len(outs) == 0 {
		return nil
	}

	query := m.sql.Insert(transactionOutputsTableName).
		Columns("tx_hash", "output_index", "value", "address", "script_hex")

	for _, out := range outs {
		query = query.Values(out.TxHash, out.Index, out.Value, out.Address, out.ScriptHex)
	}

	return m.db.ExecContext(ctx, query)
}

func (m *transactionsQ) InsertTransactionInputsBatch(ctx context.Context, ins []data.TransactionInput) error {
	if len(ins) == 0 {
		return nil
	}

	query := m.sql.Insert(transactionInputsTableName).
		Columns("tx_hash", "input_index", "prev_tx_hash", "prev_output_index")

	for _, in := range ins {
		query = query.Values(in.TxHash, in.Index, in.PrevTxHash, in.PrevIndex)
	}

	return m.db.ExecContext(ctx, query)
}

func (m *transactionsQ) SpendTransactionOutputs(ctx context.Context, ins []data.TransactionInput) error {
	if len(ins) == 0 {
		return nil
	}

	values := make([]interface{}, 0, len(ins)*3)
	placeholders := make([]string, 0, len(ins))
	argIdx := 1
	for _, in := range ins {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d)", argIdx, argIdx+1, argIdx+2))
		values = append(values, in.PrevTxHash, in.PrevIndex, in.TxHash)
		argIdx += 3
	}

	query := squirrel.Expr(
		fmt.Sprintf(
			`UPDATE %s AS t
			SET spent_by_tx_hash = v.spending_tx_hash
			FROM (VALUES %s) AS v(prev_tx_hash, prev_output_index, spending_tx_hash)
			WHERE t.tx_hash = v.prev_tx_hash 
			AND t.output_index = v.prev_output_index::integer`,
			transactionOutputsTableName,
			strings.Join(placeholders, ","),
		),
		values...,
	)

	return m.db.ExecContext(ctx, query)
}
