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
		sql: squirrel.StatementBuilder,
	}
}

type utxosQ struct {
	db  *pgdb.DB
	sql squirrel.StatementBuilderType
}

func (m *utxosQ) SelectByAddress(ctx context.Context, address string) ([]data.Utxo, error) {
	query := m.sql.Select("*").
		From(utxosTableName).
		Where("address = ?", address).
		PlaceholderFormat(squirrel.Dollar)

	var result []data.Utxo
	err := m.db.SelectContext(ctx, &result, query)
	return result, err
}

func (m *utxosQ) InsertMany(ctx context.Context, utxos []data.Utxo) ([]data.Utxo, error) {
	// TODO: Handle insertion of > 65535 / 6 ≈ 10,922 rows
	if len(utxos) == 0 {
		return nil, nil
	}

	query := m.sql.Insert(utxosTableName).
		Columns("address", "txid", "vout", "value", "block_height", "block_hash")

	for _, utxo := range utxos {
		query = query.Values(
			utxo.Address,
			utxo.TxID,
			utxo.Vout,
			utxo.Value,
			utxo.BlockHeight,
			utxo.BlockHash,
		)
	}

	query = query.Suffix("ON CONFLICT (txid, vout) DO NOTHING RETURNING *")

	var result []data.Utxo
	err := m.db.SelectContext(ctx, &result, query)
	return result, err
}
