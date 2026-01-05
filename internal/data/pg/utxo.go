package pg

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
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
		Where("spent_height IS NULL")

	var result []data.Utxo
	err := m.db.SelectContext(ctx, &result, query)
	return result, err
}

func (m *utxosQ) Spend(txHash string, txPos int, spentHeight int) error {
	query := m.sql.Update(utxosTableName).
		Set("spent_height", spentHeight).
		Where(squirrel.Eq{"tx_hash": txHash, "tx_pos": txPos})

	return m.db.Exec(query)
}

func (m *utxosQ) Insert(utxo data.Utxo) (*data.Utxo, error) {
	clauses := structs.Map(utxo)
	query := m.sql.Insert(utxosTableName).
		SetMap(clauses).
		Suffix("RETURNING *")

	var result data.Utxo
	err := m.db.Get(&result, query)
	return &result, err
}
