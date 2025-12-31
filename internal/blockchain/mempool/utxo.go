package mempool

type Utxo struct {
	TxID   string `json:"txid"`
	Vout   uint   `json:"vout"`
	Status struct {
		BlockHeight int    `json:"block_height"`
		BlockHash   string `json:"block_hash"`
	} `json:"status"`
	Value int64 `json:"value"`
}
