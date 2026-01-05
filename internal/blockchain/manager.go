package blockchain

import (
	"github.com/maphy9/btc-utxo-indexer/internal/blockchain/electrum"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
)

func NewManager(nodeAddr string, db data.MasterQ, log *logan.Entry) (*Manager, error) {
	client, err := electrum.NewClient(nodeAddr)
	if err != nil {
		return nil, err
	}

	m := &Manager{
		client: client,
		db:     db,
		log:    log,
	}

	return m, nil
}

type Manager struct {
	client *electrum.Client
	db     data.MasterQ
	log    *logan.Entry
}

func (m *Manager) Close() error {
	return m.client.Close()
}
