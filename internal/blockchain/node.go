package blockchain

import "github.com/maphy9/btc-utxo-indexer/internal/data"

type Node interface {
	GetLatestBlock() (*Block, error)

	GetAddressUtxos(address string) ([]data.Utxo, error)
}
