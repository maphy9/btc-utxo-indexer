package util

import (
	"encoding/hex"
	"errors"

	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

func ParseHeaderHex(hexStr string, height int) (*data.Header, error) {
	if len(hexStr) != 160 {
		return nil, errors.New("bad header hex")
	}
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	hash := hex.EncodeToString(Reverse(doubleHash(bytes)))
	parentHash := hex.EncodeToString(Reverse(bytes[4:36]))
	root := hex.EncodeToString(Reverse(bytes[36:68]))
	return &data.Header{
		Hash:       hash,
		ParentHash: parentHash,
		Root:       root,
		Height:     height,
	}, nil
}
