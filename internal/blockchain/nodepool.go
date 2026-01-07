package blockchain

import (
	"context"
	"errors"
	"sync"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain/electrum"
)

type NodepoolEntry struct {
	Address           string
	SSL               bool
	ReconnectAttempts uint32
}

type nodepoolEntry struct {
	NodepoolEntry
	client *electrum.Client
}

func newNodepoolEntry(entry NodepoolEntry) (nodepoolEntry, error) {
	client, err := electrum.NewClient(entry.Address, entry.SSL)
	if err != nil {
		return nodepoolEntry{}, err
	}

	return nodepoolEntry{
		NodepoolEntry: entry,
		client:        client,
	}, nil
}

func newNodepool(entries []NodepoolEntry) (*nodepool, error) {
	if len(entries) == 0 {
		return nil, errors.New("At least one node is required")
	}

	np := nodepool{
		nodes:       make([]nodepoolEntry, len(entries)),
		nodeIdx:     0,
	}
	for i, entry := range entries {
		node, err := newNodepoolEntry(entry)
		if err != nil {
			np.Close()
			return nil, err
		}
		np.nodes[i] = node
	}
	np.primaryNode = np.nodes[0]	
	
	return &np, nil
}

type nodepool struct {
	primaryNode nodepoolEntry
	nodes       []nodepoolEntry
	nodeIdx     uint64
	mu          sync.Mutex
}

func (np *nodepool) Close() error {
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

func (np *nodepool) subscribeAddress(ctx context.Context, address string) (<-chan string, error) {
	return np.primaryNode.client.SubscribeAddress(ctx, address)
}

func (np *nodepool) incrementNodeIdx() {
	np.nodeIdx = (np.nodeIdx + 1) % uint64(len(np.nodes))
}

func (np *nodepool) getPrimaryClient() (*electrum.Client, error) {
	for {
		client := np.primaryNode.client
		if client.IsHealthy() {
			return client, nil
		}
		if np.primaryNode.ReconnectAttempts < 1 {
			return nil, errors.New("Primary node is dead")
		}
		np.primaryNode.ReconnectAttempts -= 1
		reconnectedNode, err := newNodepoolEntry(np.primaryNode.NodepoolEntry)
		if err != nil {
			np.primaryNode.ReconnectAttempts = 0
			return nil, errors.New("Primary node is dead")
		}
		np.primaryNode = reconnectedNode
	}
}

func (np *nodepool) getNextClient() (*electrum.Client, error) {
	np.mu.Lock()
	defer np.mu.Unlock()

	deadCount := 0
	for {
		if deadCount == len(np.nodes) {
			return nil, errors.New("All nodes are dead")
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
		reconnectedNode, err := newNodepoolEntry(node.NodepoolEntry)
		if err != nil {
			node.ReconnectAttempts = 0
			deadCount += 1
			np.incrementNodeIdx()
			continue
		}
		np.nodes[np.nodeIdx] = reconnectedNode
		continue
	}
}

func (np *nodepool) getHeader(ctx context.Context, height int) (*electrum.Header, error) {
	client, err := np.getNextClient()
	if err != nil {
		return nil, err
	}
	return client.GetHeader(ctx, height)
}

func (np *nodepool) getTipHeight(ctx context.Context) (int, error) {
	client, err := np.getNextClient()
	if err != nil {
		return 0, err
	}
	return client.GetTipHeight(ctx)
}

func (np *nodepool) getHeaders(ctx context.Context, height, count int) ([]electrum.Header, error) {
	client, err := np.getNextClient()
	if err != nil {
		return nil, err
	}
	return client.GetHeaders(ctx, height, count)
}

func (np *nodepool) subscribeHeaders(ctx context.Context) (<-chan electrum.Header, error) {
	client, err := np.getPrimaryClient()
	if err != nil {
		return nil, err
	}
	return client.SubscribeHeaders(ctx)
}

func (np *nodepool) getTransactionMerkle(ctx context.Context, txHash string, height int) (*electrum.TransactionMerkle, error) {
	client, err := np.getNextClient()
	if err != nil {
		return nil, err
	}
	return client.GetTransactionMerkle(ctx, txHash, height)
}

func (np *nodepool) getTransaction(ctx context.Context, txHash string) (*electrum.TransactionUtxos, error) {
	client, err := np.getNextClient()
	if err != nil {
		return nil, err
	}
	return client.GetTransaction(ctx, txHash)
}

func (np *nodepool) getTransactionHeaders(ctx context.Context, address string) ([]electrum.TransactionHeader, error) {
	client, err := np.getNextClient()
	if err != nil {
		return nil, err
	}
	return client.GetTransactionHeaders(ctx, address)
}
