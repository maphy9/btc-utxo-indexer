package blockchain

import (
	"encoding/hex"
	"errors"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain/electrum"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
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

func electrumHeaderToData(hdr *electrum.Header) (*data.Header, error) {
	if len(hdr.Hex) != 160 {
		return nil, errors.New("bad header hex")
	}
	bytes, err := hex.DecodeString(hdr.Hex)
	if err != nil {
		return nil, err
	}
	hash := hex.EncodeToString(util.Reverse(util.DoubleHash(bytes)))
	parentHash := hex.EncodeToString(util.Reverse(bytes[4:36]))
	root := hex.EncodeToString(util.Reverse(bytes[36:68]))
	return &data.Header{
		Hash:       hash,
		ParentHash: parentHash,
		Root:       root,
		Height:     hdr.Height,
	}, nil
}

func electrumHeadersToData(rawHdrs []electrum.Header) ([]*data.Header, error) {
	hdrs := make([]*data.Header, len(rawHdrs))
	for i, rawHdr := range rawHdrs {
		hdr, err := electrumHeaderToData(&rawHdr)
		if err != nil {
			return nil, err
		}
		hdrs[i] = hdr
	}
	return hdrs, nil
}

func voutsToData(vouts []electrum.UtxoVout) []data.Utxo {
	utxos := make([]data.Utxo, len(vouts))
	for i, vout := range vouts {
		utxos[i] = voutToData(vout)
	}
	return utxos
}

func voutToData(vout electrum.UtxoVout) data.Utxo {
	return data.Utxo{
		Address: vout.Address,
		TxHash:  vout.TxHash,
		TxPos:   vout.N,
		Value:   vout.Value,
	}
}

func txHdrsToData(txHdrs []electrum.TransactionHeader) []data.Transaction {
	txs := make([]data.Transaction, len(txHdrs))
	for i, txHdr := range txHdrs {
		txs[i] = txHdrToData(txHdr)
	}
	return txs
}

func txHdrToData(txHdr electrum.TransactionHeader) data.Transaction {
	return data.Transaction{
		Height: txHdr.Height,
		TxHash: txHdr.TxHash,
	}
}
