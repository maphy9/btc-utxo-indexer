package pg

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const addressesTableName = "tracked_addresses"

func newAddressesQ(db *pgdb.DB) data.AddressesQ {
	return &addressesQ{
		db:  db,
		sql: squirrel.StatementBuilder,
	}
}

type addressesQ struct {
	db  *pgdb.DB
	sql squirrel.StatementBuilderType
}

func (m *addressesQ) SelectByUserID(ctx context.Context, userID int64) ([]data.Address, error) {
	query := m.sql.Select("*").
		From(addressesTableName).
		Where("user_id = ?", userID).
		PlaceholderFormat(squirrel.Dollar)

	var result []data.Address
	err := m.db.SelectContext(ctx, &result, query)
	return result, err
}

func (m *addressesQ) Insert(ctx context.Context, address data.Address) (*data.Address, error) {
	clauses := structs.Map(address)
	query := m.sql.Insert(addressesTableName).
		SetMap(clauses).
		Suffix("RETURNING *")

	var result data.Address
	err := m.db.GetContext(ctx, &result, query)
	return &result, err
}
