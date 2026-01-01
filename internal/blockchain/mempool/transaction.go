package mempool

type RawTransaction struct {
	TxID   string `json:"txid"`
	Status struct {
		BlockHeight int `json:"block_height"`
	} `json:"status"`
}

