package blockchain

type Node interface {
	GetLatestBlock() (*Block, error)
}