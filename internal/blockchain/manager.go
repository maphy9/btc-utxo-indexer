package blockchain

import (
	"log"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain/electrum"
)

func NewManager(nodeAddr string) (*Manager, error) {
	client, err := electrum.NewClient(nodeAddr)
	if err != nil {
		return nil, err
	}

	manager := &Manager{Client: client}
	return manager, nil
}

type Manager struct {
	Client *electrum.Client
}

func (m *Manager) WatchAddress(address string) error {
	notifyChan, err := m.Client.Subscribe(address)
	if err != nil {
		return err
	}

	go func() {
		for status := range notifyChan {
			log.Printf("NOTIFICATION for %s!!! status: %s", address, status)
		}
	}()

	return nil
}
