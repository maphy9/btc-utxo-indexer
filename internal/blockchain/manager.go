package blockchain

import (
	"github.com/maphy9/btc-utxo-indexer/internal/blockchain/electrum"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

func NewManager(nodeAddr string, db data.MasterQ) (*Manager, error) {
	client, err := electrum.NewClient(nodeAddr)
	if err != nil {
		return nil, err
	}

	m := &Manager{
		client: client,
		db:     db,
	}

	if err := m.subscribeToSavedAddresses(); err != nil {
		return nil, err
	}

	return m, nil
}

type Manager struct {
	client *electrum.Client
	db     data.MasterQ
}
