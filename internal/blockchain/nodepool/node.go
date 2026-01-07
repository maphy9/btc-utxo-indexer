package nodepool

import (
	"errors"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain/electrum"
)

type NodepoolEntry struct {
	Address           string
	SSL               bool
	ReconnectAttempts uint32
}

type NodeEntry struct {
	NodepoolEntry
	client *electrum.Client
}

func NewNodeEntry(entry NodepoolEntry) (*NodeEntry, error) {
	client, err := electrum.NewClient(entry.Address, entry.SSL)
	if err != nil {
		return nil, err
	}

	return &NodeEntry{
		NodepoolEntry: entry,
		client:        client,
	}, nil
}

func (np *Nodepool) getPrimaryNode() (*electrum.Client, error) {
	np.mu.Lock()
	defer np.mu.Unlock()

	for {
		client := np.primaryNode.client
		if client.IsHealthy() {
			return client, nil
		}
		if np.primaryNode.ReconnectAttempts < 1 {
			return nil, errors.New("primary node is dead")
		}
		np.primaryNode.ReconnectAttempts -= 1
		reconnectedNode, err := NewNodeEntry(np.primaryNode.NodepoolEntry)
		if err != nil {
			np.primaryNode.ReconnectAttempts = 0
			return nil, errors.New("primary node is dead")
		}
		err = np.primaryNode.client.Close()
		if err != nil {
			return nil, err
		}
		np.primaryNode = reconnectedNode
	}
}

func (np *Nodepool) getNextNode() (*electrum.Client, error) {
	np.mu.Lock()
	defer np.mu.Unlock()

	deadCount := 0
	for {
		if deadCount == len(np.nodes) {
			return nil, errors.New("all nodes are dead")
		}
		node := np.nodes[np.nodeIdx]
		if node.client.IsHealthy() {
			np.incrementNodeIdx()
			return node.client, nil
		}

		if node.ReconnectAttempts < 1 {
			deadCount += 1
			np.incrementNodeIdx()
			continue
		}
		node.ReconnectAttempts -= 1
		reconnectedNode, err := NewNodeEntry(node.NodepoolEntry)
		if err != nil {
			node.ReconnectAttempts = 0
			deadCount += 1
			np.incrementNodeIdx()
			continue
		}

		err = node.client.Close()
		if err != nil {
			return nil, err
		}
		np.nodes[np.nodeIdx] = reconnectedNode
	}
}
