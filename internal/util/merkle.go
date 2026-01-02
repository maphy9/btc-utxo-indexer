package util

import (
	"encoding/hex"
)

func VerifyMerkleProof(merkle []string, txHash string, txPos int, root string) bool {
	prevHash, _ := hex.DecodeString(txHash)
	prevHash = reverse(prevHash)
	for _, hash := range merkle {
		data, _ := hex.DecodeString(hash)
		data = reverse(data)
		if txPos%2 == 0 {
			data = append(prevHash, data...)
		} else {
			data = append(data, prevHash...)
		}
		txPos /= 2
		prevHash = doubleHash(data)
	}
	myRoot := hex.EncodeToString(reverse(prevHash))
	return root == myRoot
}
