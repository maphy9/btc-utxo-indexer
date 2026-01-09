package rpc

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"github.com/maphy9/btc-utxo-indexer/internal/util"
)

func addressToScripthash(address string) (string, error) {
	addr, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}

	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(script)
	for i, j := 0, len(hash)-1; i < j; i, j = i+1, j-1 {
		hash[i], hash[j] = hash[j], hash[i]
	}
	return hex.EncodeToString(hash[:]), nil
}

func extractTransactionData(utilTx *btcutil.Tx) *TransactionData {
	msgTx := utilTx.MsgTx()
	txHash := utilTx.Hash().String()

	tx := &TransactionData{
		Inputs:  make([]data.TransactionInput, len(msgTx.TxIn)),
		Outputs: make([]data.TransactionOutput, len(msgTx.TxOut)),
	}

	for i, in := range msgTx.TxIn {
		tx.Inputs[i] = data.TransactionInput{
			TxHash:     txHash,
			Index:      i,
			PrevTxHash: in.PreviousOutPoint.Hash.String(),
			PrevIndex:  int(in.PreviousOutPoint.Index),
		}
	}

	for i, out := range msgTx.TxOut {
		_, addrs, _, _ := txscript.ExtractPkScriptAddrs(out.PkScript, &chaincfg.MainNetParams)

		addr := ""
		if len(addrs) > 0 {
			addr = addrs[0].String()
		}

		tx.Outputs[i] = data.TransactionOutput{
			TxHash:    txHash,
			Index:     i,
			Value:     out.Value,
			Address:   addr,
			ScriptHex: hex.EncodeToString(out.PkScript),
		}
	}

	return tx
}

func txHdrToData(txHdr transactionHeader) data.Transaction {
	return data.Transaction{
		Height: txHdr.Height,
		TxHash: txHdr.TxHash,
	}
}

func headerResponseToData(hdr *headerResponse) (*data.Header, error) {
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
	rawTime := binary.LittleEndian.Uint32(bytes[68:72])
	time := time.Unix(int64(rawTime), 0)
	return &data.Header{
		Hash:       hash,
		ParentHash: parentHash,
		Root:       root,
		Height:     hdr.Height,
		CreatedAt:  time,
	}, nil
}
