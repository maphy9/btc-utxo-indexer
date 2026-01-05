package blockchain

func (m *Manager) SubscribeAddress(address string) error {
	notifyChan, err := m.client.SubscribeAddress(address)
	if err != nil {
		m.log.WithError(err).Errorf("failed to subscribe to address (%s)", address)
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
			m.log.WithError(err).Errorf("failed to get address status (%s)", address)
			continue
		}
		if oldStatus == status {
			continue
		}
		err = m.syncHistory(address)
		if err != nil {
			m.log.WithError(err).Errorf("failed to sync address (%s)", address)
			continue
		}
		m.log.Infof("finished sync for address (%s)", address)
		err = m.db.Addresses().UpdateStatus(address, status)
		if err != nil {
			m.log.WithError(err).Errorf("failed to update address status (%s)", address)
			continue
		}
	}
}
