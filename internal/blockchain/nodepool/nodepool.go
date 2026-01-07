package nodepool

import (
	"errors"
	"sync"
)

func NewNodepool(entries []NodepoolEntry) (*Nodepool, error) {
	if len(entries) == 0 {
		return nil, errors.New("At least one node is required")
	}

	np := Nodepool{
		nodes:   make([]*NodeEntry, len(entries)),
		nodeIdx: 0,
	}
	for i, entry := range entries {
		node, err := NewNodeEntry(entry)
		if err != nil {
			np.Close()
			return nil, err
		}
		np.nodes[i] = node
	}
	np.primaryNode = np.nodes[0]

	return &np, nil
}

type Nodepool struct {
	primaryNode *NodeEntry
	nodes       []*NodeEntry
	nodeIdx     uint64
	mu          sync.Mutex
}

func (np *Nodepool) incrementNodeIdx() {
	np.nodeIdx = (np.nodeIdx + 1) % uint64(len(np.nodes))
}

func (np *Nodepool) Close() error {
	for _, node := range np.nodes {
		if node.client == nil {
			continue
		}
		err := node.client.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
