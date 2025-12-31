package blockchain

type Node interface {
	GetLatestBlock() (*Block, error)

	GetAddressUtxos(address string) ([]Utxo, error)
}