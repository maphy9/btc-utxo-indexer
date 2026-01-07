package nodepool

type HealthStatus struct {
	NodeAddress string
	IsHealthy   bool
}

func (np *Nodepool) GetHealthStatuses() []HealthStatus {
	np.mu.Lock()
	defer np.mu.Unlock()
	healthStatuses := make([]HealthStatus, len(np.nodes))
	for i, node := range np.nodes {
		healthStatuses[i] = HealthStatus{
			NodeAddress: node.Address,
			IsHealthy:   node.client.IsHealthy(),
		}
	}
	return healthStatuses
}
