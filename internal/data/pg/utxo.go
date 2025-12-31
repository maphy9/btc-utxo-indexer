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
		db: db,
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

func (m *utxosQ) Insert(ctx context.Context, utxo data.Utxo) (*data.Utxo, error) {
	clauses := structs.Map(utxo)
	query := m.sql.Insert(utxosTableName).
		SetMap(clauses).
		Suffix("RETURNING *")

	var result data.Utxo
	err := m.db.GetContext(ctx, &result, query)
	return &result, err
}