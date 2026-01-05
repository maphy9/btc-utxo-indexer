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
		sql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

type addressesQ struct {
	db  *pgdb.DB
	sql squirrel.StatementBuilderType
}

func (m *addressesQ) GetUserAddresses(ctx context.Context, userID int64) ([]data.Address, error) {
	query := m.sql.Select("a.*").
		From(addressesTableName+" a").
		Join(userAddressesTableName+" ua ON a.id = ua.address_id").
		Where("ua.user_id = ?", userID)

	var result []data.Address
	err := m.db.SelectContext(ctx, &result, query)
	return result, err
}

func (m *addressesQ) GetUserAddress(ctx context.Context, userID int64, address string) (*data.UserAddress, error) {
	query := m.sql.Select("ua.*").
		From(addressesTableName+" a").
		Join(userAddressesTableName+" ua ON a.id = ua.address_id").
		Where("user_id = ?", userID).
		Where("a.address = ?", address)

	var result data.UserAddress
	err := m.db.GetContext(ctx, &result, query)
	return &result, err
}

func (m *addressesQ) GetAllAddresses() ([]string, error) {
	query := m.sql.Select("address").
		From(addressesTableName)

	var result []string
	err := m.db.Select(&result, query)
	return result, err
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

func (m *addressesQ) GetStatus(address string) (string, error) {
	query := m.sql.Select("status").
		From(addressesTableName).
		Where("address = ?", address)
	var status string
	err := m.db.Get(&status, query)
	return status, err
}

func (m *addressesQ) GetBalance(ctx context.Context, address string) (int64, error) {
	query := m.sql.Select("SUM(value)").
		From(utxosTableName).
		Where("address = ?", address).
		Where("spent_height IS NULL")

	var result int64
	err := m.db.GetContext(ctx, &result, query)
	return result, err
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

func (m *addressesQ) UpdateStatus(address, status string) error {
	query := m.sql.Update(addressesTableName).
		Set("status", status).
		Where("address = ?", address)
	return m.db.Exec(query)
}

func (m *addressesQ) Exists(address string) (bool, error) {
	query := m.sql.Select("COUNT(*)").
		From(addressesTableName).
		Where("address = ?", address)

	var result int
	err := m.db.Get(&result, query)
	return result > 0, err
}
