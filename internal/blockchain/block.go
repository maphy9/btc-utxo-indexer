package blockchain

type Block struct {
	Hash       string
	Height     int
	ParentHash string
	MerkleRoot string
}
