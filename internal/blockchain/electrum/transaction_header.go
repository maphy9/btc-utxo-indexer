package electrum

import (
	"encoding/json"
)

type TransactionHeader struct {
	Height int    `json:"height"`
	TxHash string `json:"tx_hash"`
}

func (c *Client) GetTransactionHeaders(address string) ([]TransactionHeader, error) {
	scripthash, err := addressToScripthash(address)
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
