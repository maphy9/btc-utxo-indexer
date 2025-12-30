package blockchain

type Block struct {
	Hash string `json:"id"`
	Height int `json:"height"`
	PreviousBlockHash string `json:"previousblockhash"`
	MerkleRoot string `json:"merkle_root"`
}