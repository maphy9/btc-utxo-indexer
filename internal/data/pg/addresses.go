package pg

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	addressesTableName     = "addresses"
	userAddressesTableName = "user_addresses"
)

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

func (m *addressesQ) GetUserAddresses(ctx context.Context, userID int64) ([]data.Address, error) {
	query := m.sql.Select("a.*").
		From(addressesTableName + " a").
		Join(userAddressesTableName + " ua ON a.id = ua.address_id").
		Where("ua.user_id = ?", userID).
		PlaceholderFormat(squirrel.Dollar)

	var result []data.Address
	err := m.db.SelectContext(ctx, &result, query)
	return result, err
}

func (m *addressesQ) GetUserAddress(ctx context.Context, userID int64, address string) (*data.UserAddress, error) {
	query := m.sql.Select("ua.*").
		From(addressesTableName+" a").
		Join(userAddressesTableName+" ua ON a.id = ua.address_id").
		Where("user_id = ?", userID).
		Where("a.address = ?", address).
		PlaceholderFormat(squirrel.Dollar)

	var result data.UserAddress
	err := m.db.GetContext(ctx, &result, query)
	return &result, err
}

func (m *addressesQ) InsertAddress(ctx context.Context, address string) (*data.Address, error) {
	query := m.sql.Insert(addressesTableName).
		Columns("address").
		Values(address).
		Suffix(`
			ON CONFLICT (address) DO
			UPDATE SET address = EXCLUDED.address
			RETURNING *
		`)

	var result data.Address
	err := m.db.GetContext(ctx, &result, query)
	return &result, err
}

func (m *addressesQ) InsertUserAddress(ctx context.Context, userAddress data.UserAddress) (*data.UserAddress, error) {
	clauses := structs.Map(userAddress)
	query := m.sql.Insert(userAddressesTableName).
		SetMap(clauses).
		Suffix("RETURNING *")

	var result data.UserAddress
	err := m.db.GetContext(ctx, &result, query)
	return &result, err
}
