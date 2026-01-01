package mempool

import "github.com/maphy9/btc-utxo-indexer/internal/data"

type RawUtxo struct {
	TxID   string `json:"txid"`
	Vout   uint   `json:"vout"`
	Status struct {
		BlockHeight int `json:"block_height"`
	} `json:"status"`
	Value int64 `json:"value"`
}

func mapRawUtxo(utxo RawUtxo, address string) data.Utxo {
	return data.Utxo{
		Address:     address,
		TxID:        utxo.TxID,
		Vout:        utxo.Vout,
		Value:       utxo.Value,
		BlockHeight: utxo.Status.BlockHeight,
	}
}

func mapRawUtxos(utxos []RawUtxo, address string) []data.Utxo {
	mappedUtxos := make([]data.Utxo, 0, len(utxos))
	for _, utxo := range utxos {
		mappedUtxos = append(mappedUtxos, mapRawUtxo(utxo, address))
	}
	return mappedUtxos
}