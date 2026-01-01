package electrum

import (
	"encoding/json"

	"github.com/maphy9/btc-utxo-indexer/internal/util"
)

type Utxo struct {
	Height int    `json:"height"`
	TxPos  int    `json:"tx_pos"`
	TxHash string `json:"tx_hash"`
	Value  int64  `json:"value"`
}

func (c *Client) GetUtxos(address string) ([]Utxo, error) {
	scripthash, err := util.AddressToScripthash(address)
	if err != nil {
		return nil, err
	}

	rawUtxos, err := c.request("blockchain.scripthash.listunspent", []any{scripthash})
	if err != nil {
		return nil, err
	}

	var utxos []Utxo
	err = json.Unmarshal(rawUtxos, &utxos)
	if err != nil {
		return nil, err
	}
	return utxos, nil
}
