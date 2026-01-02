package blockchain

import "log"

func (m *Manager) subscribeSavedAddresses() error {
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

func (m *Manager) SubscribeAddress(address string) error {
	notifyChan, err := m.client.SubscribeAddress(address)
	if err != nil {
		log.Printf("Failed to subscribe to an address: %v", err)
		return err
	}

	go m.watchAddress(address, notifyChan)
	return nil
}

func (m *Manager) watchAddress(address string, notifyChan <-chan string) {
	for status := range notifyChan {
		oldStatus, err := m.db.Addresses().UpdateStatus(address, status)
		if err != nil {
			log.Printf("Failed to get address status: %v", err)
			continue
		}
		if oldStatus == status {
			continue
		}
		m.synchronizeHistory(address)
	}
}
