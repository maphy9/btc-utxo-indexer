package pg

import (
	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

func NewMasterQ(db *pgdb.DB) data.MasterQ {
	return &masterQ{db}
}

type masterQ struct {
	db *pgdb.DB
}

func (m *masterQ) Users() data.UsersQ {
	return newUsersQ(m.db)
}

func (m *masterQ) Addresses() data.AddressesQ {
	return newAddressesQ(m.db)
}

func (m *masterQ) Utxos() data.UtxosQ {
	return newUtxosQ(m.db)
}

func (m *masterQ) Transaction(fn func(q data.MasterQ) error) error {
	return m.db.Transaction(func() error {
		return fn(m)
	})
}
