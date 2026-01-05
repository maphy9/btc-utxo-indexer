package blockchain

import (
	"context"
	"sync"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain/electrum"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
)

func NewManager(client *electrum.Client, db data.MasterQ, log *logan.Entry) (*Manager, error) {
	ctx, cancel := context.WithCancel(context.Background())

	m := &Manager{
		ctx:    ctx,
		cancel: cancel,
		client: client,
		db:     db,
		log:    log,
	}

	return m, nil
}

type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	client *electrum.Client
	db     data.MasterQ
	log    *logan.Entry
}

func (m *Manager) Close() error {
	m.cancel()
	m.wg.Wait()
	return m.client.Close()
}
