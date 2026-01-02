package electrum

import (
	"encoding/json"

	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

type UtxoVout struct {
	Value        float64 `json:"value"`
	N            int     `json:"n"`
	ScriptPubKey struct {
		Addresses []string `json:"addresses"`
	} `json:"scriptPubKey"`
}

func (utxoVout *UtxoVout) ToData(txHash string, height int) data.Utxo {
	sats := int64(utxoVout.Value * 100_000_000)
	return data.Utxo{
		Address:       utxoVout.ScriptPubKey.Addresses[0],
		TxHash:        txHash,
		TxPos:         utxoVout.N,
		Value:         sats,
		CreatedHeight: height,
	}
}

type Transaction struct {
	TxID string `json:"txid"`
	Vin  []struct {
		TxID string `json:"txid"`
		Vout int    `json:"vout"`
	} `json:"vin"`
	Vout []UtxoVout `json:"vout"`
}

type TransactionMerkle struct {
	Height int      `json:"block_height"`
	Merkle []string `json:"merkle"`
	Pos    int      `json:"pos"`
}

func (c *Client) GetTransaction(txHash string) (*Transaction, error) {
	rawTx, err := c.request("blockchain.transaction.get", []any{txHash, true})
	if err != nil {
		return nil, err
	}

	var tx Transaction
	err = json.Unmarshal(rawTx, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (c *Client) GetTransactionMerkle(txHash string, height int) (*TransactionMerkle, error) {
	rawTxMerkle, err := c.request("blockchain.transaction.get_merkle", []any{txHash, height})
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
