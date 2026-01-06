package blockchain

import "context"

func (m *Manager) SubscribeAddress(address string) error {
	notifyChan, err := m.client.SubscribeAddress(m.ctx, address)
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

func (m *Manager) processAddress(ctx context.Context, address, status string) error {
	oldStatus, err := m.db.Addresses().GetStatus(address)
	if err != nil {
		return err
	}

	if oldStatus == status {
		return err
	}

	err = m.syncHistory(m.ctx, address)
	if err != nil {
		return err
	}

	err = m.db.Addresses().UpdateStatus(address, status)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) watchAddress(address string, notifyChan <-chan string) {
	m.wg.Add(1)
	defer m.wg.Done()

	for {
		select {
		case <-m.ctx.Done():
			return
		case status, ok := <-notifyChan:
			if !ok {
				m.log.Info("Addresses channed was closed")
				return
			}
			err := m.processAddress(m.ctx, address, status)
			if err != nil {
				m.log.WithError(err).Errorf("failed to process addess status update (%s)", address)
				continue
			}
			m.log.Infof("processed address (%s)", address)
		}
	}

}
