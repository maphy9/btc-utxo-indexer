package electrum

import (
	"context"
	"encoding/hex"
	"encoding/json"

	"github.com/btcsuite/btcd/btcutil"
)

type UtxoVout struct {
	Value        float64
	N            int
	ScriptPubKey struct {
		Addresses []string
	}
}

type Transaction struct {
	TxID string
	Vin  []struct {
		TxID string
		Vout int
	}
	Vout []UtxoVout
}

type TransactionMerkle struct {
	Height int      `json:"block_height"`
	Merkle []string `json:"merkle"`
	Pos    int      `json:"pos"`
}

func (c *Client) GetTransaction(ctx context.Context, txHash string) (*Transaction, error) {
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
	return btcutilToTransaction(btcutilTx), nil
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