package data

type TransactionsQ interface {
	Exists(txHash string) (bool, error)
	Insert(tx Transaction) (*Transaction, error)
}

type Transaction struct {
	TxHash string `db:"tx_hash" structs:"tx_hash" json:"tx_hash"`
	Height int    `db:"height" structs:"height" json:"height"`
}
