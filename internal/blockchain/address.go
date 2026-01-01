package blockchain

import "log"

func (m *Manager) SubscribeAddress(address string) error {
	notifyChan, err := m.client.Subscribe(address)
	if err != nil {
		log.Fatal(err)
		return err
	}

	go m.watchAddress(address, notifyChan)
	return nil
}

func (m *Manager) subscribeToSavedAddresses() error {
	addresses, err := m.db.Addresses().GetAllAddresses()
	if err != nil {
		return err
	}

	for _, address := range addresses {
		if err := m.SubscribeAddress(address); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) watchAddress(address string, notifyChan <-chan string) {
	for status := range notifyChan {
		oldStatus, err := m.db.Addresses().UpdateStatus(address, status)
		if err != nil {
			log.Printf("Error while syncing address status: %v", err)
			continue
		}
		if oldStatus == status {
			continue
		}
		m.synchronizeHistory(address)
	}
}
