package nodepool

import (
	"context"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain/rpc"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

func (np *Nodepool) SubscribeAddress(ctx context.Context, address string) (<-chan string, error) {
	client, err := np.getPrimaryNode()
	if err != nil {
		return nil, err
	}
	return client.SubscribeAddress(ctx, address)
}

func (np *Nodepool) SubscribeHeaders(ctx context.Context) (<-chan *data.Header, error) {
	client, err := np.getPrimaryNode()
	if err != nil {
		return nil, err
	}
	return client.SubscribeHeaders(ctx)
}

func (np *Nodepool) GetHeader(ctx context.Context, height int) (*data.Header, error) {
	client, err := np.getNextNode()
	if err != nil {
		return nil, err
	}
	return client.GetHeader(ctx, height)
}

func (np *Nodepool) GetTipHeight(ctx context.Context) (int, error) {
	client, err := np.getNextNode()
	if err != nil {
		return 0, err
	}
	return client.GetTipHeight(ctx)
}

func (np *Nodepool) GetHeaders(ctx context.Context, height, count int) ([]*data.Header, error) {
	client, err := np.getNextNode()
	if err != nil {
		return nil, err
	}
	return client.GetHeaders(ctx, height, count)
}

func (np *Nodepool) GetTransactionMerkle(ctx context.Context, txHash string, height int) (*rpc.TransactionMerkle, error) {
	client, err := np.getNextNode()
	if err != nil {
		return nil, err
	}
	return client.GetTransactionMerkle(ctx, txHash, height)
}

func (np *Nodepool) GetTransactionData(ctx context.Context, txHash string) (*rpc.TransactionData, error) {
	client, err := np.getNextNode()
	if err != nil {
		return nil, err
	}
	return client.GetTransactionData(ctx, txHash)
}

func (np *Nodepool) GetTransactionHeaders(ctx context.Context, address string) ([]data.Transaction, error) {
	client, err := np.getNextNode()
	if err != nil {
		return nil, err
	}
	return client.GetTransactionHeaders(ctx, address)
}
