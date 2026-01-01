package blockchain

import (
	"log"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain/electrum"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

func NewManager(nodeAddr string, db data.MasterQ) (*Manager, error) {
	client, err := electrum.NewClient(nodeAddr)
	if err != nil {
		return nil, err
	}

	manager := &Manager{
		client: client,
		db:     db,
	}
	return manager, nil
}

type Manager struct {
	client *electrum.Client
	db     data.MasterQ
}

func (m *Manager) SubscribeAddress(address string) error {
	notifyChan, err := m.client.Subscribe(address)
	if err != nil {
		log.Fatal(err)
		return err
	}

	go m.watchAddress(address, notifyChan)
	return nil
}

func (m *Manager) watchAddress(address string, notifyChan <-chan string) {
	for status := range notifyChan {
		log.Printf("NOTIFICATION for %s!!! status: %s", address, status)

		oldStatus, err := m.db.Addresses().UpdateStatus(address, status)
		if err != nil {
			log.Printf("Error while syncing address status: %v", err)
			continue
		}
		if oldStatus != status {
			log.Printf("New status for %s", address)
		}
	}
}
