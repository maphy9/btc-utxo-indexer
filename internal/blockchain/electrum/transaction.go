package electrum

import (
	"encoding/json"

	"github.com/maphy9/btc-utxo-indexer/internal/util"
)

type TransactionHeader struct {
	Height int    `json:"height"`
	TxHash string `json:"tx_hash"`
}

type Transaction struct {
	TxID string `json:"txid"`
	Vin  []struct {
		TxID string `json:"txid"`
		Vout int    `json:"vout"`
	} `json:"vin"`
	Vout []struct {
		Value        float64 `json:"value"`
		N            int     `json:"n"`
		ScriptPubKey struct {
			Addresses []string `json:"addresses"`
		} `json:"scriptPubKey"`
	} `json:"vout"`
}

type TransactionMerkle struct {
	Height int      `json:"block_height"`
	Merkle []string `json:"merkle"`
	Pos    int      `json:"pos"`
}

func (c *Client) GetTransactionHeaders(address string) ([]TransactionHeader, error) {
	scripthash, err := util.AddressToScripthash(address)
	if err != nil {
		return nil, err
	}

	rawTxHdrs, err := c.request("blockchain.scripthash.get_history", []any{scripthash})
	if err != nil {
		return nil, err
	}

	var txHdrs []TransactionHeader
	err = json.Unmarshal(rawTxHdrs, &txHdrs)
	if err != nil {
		return nil, err
	}
	return txHdrs, nil
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
