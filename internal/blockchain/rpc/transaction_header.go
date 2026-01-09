package rpc

import (
	"context"
	"encoding/json"

	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

type transactionHeader struct {
	Height int    `json:"height"`
	TxHash string `json:"tx_hash"`
}

func (c *Client) GetTransactionHeaders(ctx context.Context, address string) ([]data.Transaction, error) {
	scripthash, err := addressToScripthash(address)
	if err != nil {
		return nil, err
	}

	rawTxHdrs, err := c.request(ctx, "blockchain.scripthash.get_history", []any{scripthash})
	if err != nil {
		return nil, err
	}

	var txHdrs []transactionHeader
	err = json.Unmarshal(rawTxHdrs, &txHdrs)
	if err != nil {
		return nil, err
	}

	txHdrsFiltered := make([]data.Transaction, 0, len(txHdrs))
	for _, txHdr := range txHdrs {
		if txHdr.Height == -1 {
			continue // TODO: handle mempool headers
		}
		txHdrsFiltered = append(txHdrsFiltered, txHdrToData(txHdr))
	}
	return txHdrsFiltered, nil
}
