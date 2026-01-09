package rpc

import (
	"context"
	"encoding/hex"
	"encoding/json"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

type TransactionData struct {
	Inputs  []data.TransactionInput
	Outputs []data.TransactionOutput
}

type TransactionMerkle struct {
	Height int      `json:"block_height"`
	Merkle []string `json:"merkle"`
	Pos    int      `json:"pos"`
}

func (c *Client) GetTransactionData(ctx context.Context, txHash string) (*TransactionData, error) {
	rawRes, err := c.request(ctx, "blockchain.transaction.get", []any{txHash})
	if err != nil {
		return nil, err
	}

	var hexStr string
	if err := json.Unmarshal(rawRes, &hexStr); err != nil {
		return nil, err
	}
	txBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	btcutilTx, err := btcutil.NewTxFromBytes(txBytes)
	if err != nil {
		return nil, err
	}
	return extractTransactionData(btcutilTx), nil
}

func (c *Client) GetTransactionMerkle(ctx context.Context, txHash string, height int) (*TransactionMerkle, error) {
	rawTxMerkle, err := c.request(ctx, "blockchain.transaction.get_merkle", []any{txHash, height})
	if err != nil {
		return nil, err
	}

	var txMerkle TransactionMerkle
	err = json.Unmarshal(rawTxMerkle, &txMerkle)
	if err != nil {
		return nil, err
	}
	return &txMerkle, nil
}
