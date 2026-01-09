package blockchain

import (
	"encoding/hex"

	"github.com/maphy9/btc-utxo-indexer/internal/util"
)

func verifyMerkleProof(merkle []string, txHash string, txPos int, root string) bool {
	prevHash, _ := hex.DecodeString(txHash)
	prevHash = util.Reverse(prevHash)
	for _, hash := range merkle {
		data, _ := hex.DecodeString(hash)
		data = util.Reverse(data)
		if txPos%2 == 0 {
			data = append(prevHash, data...)
		} else {
			data = append(data, prevHash...)
		}
		txPos /= 2
		prevHash = util.DoubleHash(data)
	}
	myRoot := hex.EncodeToString(util.Reverse(prevHash))
	return root == myRoot
}
