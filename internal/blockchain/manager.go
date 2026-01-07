package blockchain

import (
	"context"
	"sync"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain/nodepool"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
)

func NewManager(entries []nodepool.NodepoolEntry, db data.MasterQ, log *logan.Entry) (*Manager, error) {
	np, err := nodepool.NewNodepool(entries)
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
	np     *nodepool.Nodepool
	db     data.MasterQ
	log    *logan.Entry
}

func (m *Manager) Close() error {
	m.cancel()
	m.wg.Wait()
	return m.np.Close()
}

func (m *Manager) LogNodesHealth() {
	healthStatuses := m.np.GetHealthStatuses()
	for _, healthStatus := range healthStatuses {
		m.log.Infof("Node %s: isHealthy: %t", healthStatus.NodeAddress, healthStatus.IsHealthy)
	}
}