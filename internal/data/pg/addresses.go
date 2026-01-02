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

func (m *addressesQ) InsertUserAddress(ctx context.Context, userAddress data.UserAddress) (*data.UserAddress, error) {
	clauses := structs.Map(userAddress)
	query := m.sql.Insert(userAddressesTableName).
		SetMap(clauses).
		Suffix("RETURNING *")

	var result data.UserAddress
	err := m.db.GetContext(ctx, &result, query)
	return &result, err
}

func (m *addressesQ) UpdateStatus(address, status string) (string, error) {
	var oldStatus string
	err := m.db.Transaction(func() error {
		stmt := m.sql.Select("status").
			From(addressesTableName).
			Where(squirrel.Eq{"address": address}).
			Suffix("FOR UPDATE")

		if err := m.db.Get(&oldStatus, stmt); err != nil {
			return err
		}

		update := m.sql.Update(addressesTableName).
			Set("status", status).
			Where(squirrel.Eq{"address": address})

		err := m.db.Exec(update)
		return err
	})
	return oldStatus, err
}

func (m *addressesQ) Exists(address string) (bool, error) {
	query := m.sql.Select("COUNT(*)").
		From(addressesTableName).
		Where("address = ?")

	var result int
	err := m.db.Get(&result, query)
	return result > 0, err
}
