package util

import "crypto/sha256"

func doubleHash(data []byte) []byte {
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:]
}

func reverse(data []byte) []byte {
	res := make([]byte, len(data))
	for i, b := range data {
		res[len(data)-i-1] = b
	}
	return res
}
