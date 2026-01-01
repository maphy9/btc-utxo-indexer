package mempool

import "time"

type RawBlock struct {
	Height int `json:"height"`
	Hash string `json:"id"`
	ParentHash string `json:"previousblockhash"`
	Timestamp time.Time `json:"timestamp"`
}