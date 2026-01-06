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

func btcutilToTransaction(utilTx *btcutil.Tx) *Transaction {
	msgTx := utilTx.MsgTx()

	tx := &Transaction{
		TxID: utilTx.Hash().String(),
		Vin: make([]struct {
			TxID string
			Vout int
		}, len(msgTx.TxIn)),
		Vout: make([]UtxoVout, len(msgTx.TxOut)),
	}

	for i, in := range msgTx.TxIn {
		tx.Vin[i].TxID = in.PreviousOutPoint.Hash.String()
		tx.Vin[i].Vout = int(in.PreviousOutPoint.Index)
	}

	for i, out := range msgTx.TxOut {
		valBTC := float64(out.Value) / 100_000_000.0
		_, addrs, _, _ := txscript.ExtractPkScriptAddrs(out.PkScript, &chaincfg.MainNetParams)

		var addrStrings []string
		for _, addr := range addrs {
			addrStrings = append(addrStrings, addr.EncodeAddress())
		}

		tx.Vout[i] = UtxoVout{
			Value: valBTC,
			N:     i,
			ScriptPubKey: struct {
				Addresses []string
			}{
				Addresses: addrStrings,
			},
		}
	}

	return tx
}
