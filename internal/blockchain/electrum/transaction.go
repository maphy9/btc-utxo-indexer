package electrum

import (
	"context"
	"encoding/hex"
	"encoding/json"

	"github.com/btcsuite/btcd/btcutil"
)

type UtxoVout struct {
	TxHash string
	Value        int64
	N            int
	Address string
}

type UtxoVin struct {
	SpentTxHash string
	TxHash string
	Vout int
}

type TransactionUtxos struct {
	Vins  []UtxoVin
	Vouts []UtxoVout
}

type TransactionMerkle struct {
	Height int      `json:"block_height"`
	Merkle []string `json:"merkle"`
	Pos    int      `json:"pos"`
}

func (c *Client) GetTransaction(ctx context.Context, txHash string) (*TransactionUtxos, error) {
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
	return extractTransactionUtxos(btcutilTx), nil
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
