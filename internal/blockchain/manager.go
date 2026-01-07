package blockchain

import (
	"context"
	"sync"

	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
)

func NewManager(entries []NodepoolEntry, db data.MasterQ, log *logan.Entry) (*Manager, error) {
	np, err := newNodepool(entries)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	m := &Manager{
		ctx:    ctx,
		cancel: cancel,
		np:     np,
		db:     db,
		log:    log,
	}

	return m, nil
}

type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	np     *nodepool
	db     data.MasterQ
	log    *logan.Entry
}

func (m *Manager) Close() error {
	m.cancel()
	m.wg.Wait()
	return m.np.Close()
}
