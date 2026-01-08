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

func (np *Nodepool) GetHealthyCount() int {
	np.mu.Lock()
	defer np.mu.Unlock()
	count := 0
	for _, node := range np.nodes {
		if node.client.IsHealthy() {
			count += 1
		}
	}
	return count
}

func (np *Nodepool) incrementNodeIdx() {
	np.nodeIdx = (np.nodeIdx + 1) % uint64(len(np.nodes))
}
