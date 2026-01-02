package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func doubleHash(data []byte) []byte {
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:]
}

func reverse(data []byte) []byte {
	res := make([]byte, len(data))
	for i, b := range data {
		res[len(data) - i - 1] = b
	}
	return res
}

func VerifyMerkleProof(merkle []string, txHash string, txPos int, root string) bool {
	prevHash, _ := hex.DecodeString(txHash)
	prevHash = reverse(prevHash)
	for _, hash := range merkle {
		data, _ := hex.DecodeString(hash)
		data = reverse(data)
		if txPos % 2 == 0 {
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