package responses

import "github.com/maphy9/btc-utxo-indexer/internal/data"

type GetUtxosResponse struct {
	Utxos []data.Utxo `json:"utxos"`
}

func NewGetUtxosResponse(utxos []data.Utxo) GetUtxosResponse {
	return GetUtxosResponse{utxos}
}