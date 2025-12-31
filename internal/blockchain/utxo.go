package blockchain

type Utxo struct {
	TxID        string
	Vout        uint
	Value       int64
	BlockHeight int
	BlockHash   string
}
