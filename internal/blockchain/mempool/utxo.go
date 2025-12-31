package mempool

type RawUtxo struct {
	TxID   string `json:"txid"`
	Vout   uint   `json:"vout"`
	Status struct {
		BlockHeight int    `json:"block_height"`
	} `json:"status"`
	Value int64 `json:"value"`
}
