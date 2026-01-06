package electrum

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
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

func extractTransactionUtxos(utilTx *btcutil.Tx) *TransactionUtxos {
	msgTx := utilTx.MsgTx()
	txHash := utilTx.Hash().String()

	tx := &TransactionUtxos{
		Vins:  make([]UtxoVin, len(msgTx.TxIn)),
		Vouts: make([]UtxoVout, len(msgTx.TxOut)),
	}

	for i, in := range msgTx.TxIn {
		tx.Vins[i].TxHash = in.PreviousOutPoint.Hash.String()
		tx.Vins[i].Vout = int(in.PreviousOutPoint.Index)
		tx.Vins[i].SpentTxHash = txHash
	}

	for i, out := range msgTx.TxOut {
		_, addrs, _, _ := txscript.ExtractPkScriptAddrs(out.PkScript, &chaincfg.MainNetParams)

		addr := ""
		if len(addrs) > 0 {
			addr = addrs[0].String()
		}

		tx.Vouts[i] = UtxoVout{
			TxHash:  txHash,
			Value:   out.Value,
			N:       i,
			Address: addr,
		}
	}

	return tx
}
