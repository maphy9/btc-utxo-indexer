package electrum

import (
	"context"
	"encoding/json"
)

type TransactionHeader struct {
	Height int    `json:"height"`
	TxHash string `json:"tx_hash"`
}

func (c *Client) GetTransactionHeaders(ctx context.Context, address string) ([]TransactionHeader, error) {
	scripthash, err := addressToScripthash(address)
	if err != nil {
		return nil, err
	}

	rawTxHdrs, err := c.request(ctx, "blockchain.scripthash.get_history", []any{scripthash})
	if err != nil {
		return nil, err
	}

	var txHdrs []TransactionHeader
	err = json.Unmarshal(rawTxHdrs, &txHdrs)
	if err != nil {
		return nil, err
	}

	txHdrsFiltered := make([]TransactionHeader, 0, len(txHdrs))
	for _, txHdr := range txHdrs {
		if txHdr.Height == -1 {
			continue  // TODO: handle mempool headers
		}
		txHdrsFiltered = append(txHdrsFiltered, txHdr)
	}
	return txHdrsFiltered, nil
}
