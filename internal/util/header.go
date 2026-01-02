package util

import (
	"encoding/hex"
	"errors"

	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

func ParseHeaderHex(hexStr string, height int) (data.Header, error) {
	if len(hexStr) != 160 {
		return data.Header{}, errors.New("Bad header hex")
	}
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return data.Header{}, err
	}
	hash := hex.EncodeToString(reverse(doubleHash(bytes)))
	parentHash := hex.EncodeToString(reverse(bytes[4:36]))
	root := hex.EncodeToString(reverse(bytes[36:68]))
	return data.Header{
		Hash:       hash,
		ParentHash: parentHash,
		Root:       root,
		Height:     height,
	}, nil
}
