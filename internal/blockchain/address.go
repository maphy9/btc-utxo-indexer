package blockchain

import "log"

func (m *Manager) SubscribeAddress(address string) error {
	notifyChan, err := m.client.SubscribeAddress(address)
	if err != nil {
		log.Printf("Failed to subscribe to an address: %v", err)
		return err
	}

	go m.watchAddress(address, notifyChan)
	return nil
}

func (m *Manager) SubscribeSavedAddresses() error {
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
		oldStatus, err := m.db.Addresses().GetStatus(address)
		if err != nil {
			log.Printf("Failed to get address status: %v", err)
			continue
		}
		if oldStatus == status {
			continue
		}
		err = m.syncHistory(address)
		if err != nil {
			log.Printf("Failed to sync address: %v", err)
			continue
		}
		log.Printf("Finished sync for address %s", address)
		err = m.db.Addresses().UpdateStatus(address, status)
		if err != nil {
			log.Printf("Failed to update address status: %v", err)
			continue
		}
	}
}
