package pg

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const utxosTableName = "utxos"

func newUtxosQ(db *pgdb.DB) data.UtxosQ {
	return &utxosQ{
		db:  db,
		sql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

type utxosQ struct {
	db  *pgdb.DB
	sql squirrel.StatementBuilderType
}

func (m *utxosQ) GetActiveByAddress(ctx context.Context, address string) ([]data.Utxo, error) {
	query := m.sql.Select("*").
		From(utxosTableName).
		Where("address = ?", address).
		Where("spent_tx_hash IS NULL")

	var result []data.Utxo
	err := m.db.SelectContext(ctx, &result, query)
	return result, err
}

func (m *utxosQ) Spend(ctx context.Context, txHash string, txPos int, spentTxHash string) error {
	query := m.sql.Update(utxosTableName).
		Set("spent_tx_hash", spentTxHash).
		Where("tx_hash = ? AND tx_pos = ?", txHash, txPos)

	return m.db.ExecContext(ctx, query)
}

func (m *utxosQ) InsertBatch(ctx context.Context, utxos []data.Utxo) error {
	if len(utxos) == 0 {
		return nil
	}

	query := m.sql.Insert(utxosTableName).
		Columns("address", "tx_hash", "tx_pos", "value").
		Suffix("ON CONFLICT DO NOTHING")

	for _, utxo := range utxos {
		query = query.Values(utxo.Address, utxo.TxHash, utxo.TxPos, utxo.Value)
	}

	return m.db.ExecContext(ctx, query)
}
