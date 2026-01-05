package electrum

import (
	"encoding/hex"
	"encoding/json"
	"log"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

type UtxoVout struct {
	Value        float64 `json:"value"`
	N            int     `json:"n"`
	ScriptPubKey struct {
		Addresses []string `json:"addresses"`
	} `json:"scriptPubKey"`
}

type Transaction struct {
	TxID string `json:"txid"`
	Vin  []struct {
		TxID string `json:"txid"`
		Vout int    `json:"vout"`
	} `json:"vin"`
	Vout []UtxoVout `json:"vout"`
}

type TransactionMerkle struct {
	Height int      `json:"block_height"`
	Merkle []string `json:"merkle"`
	Pos    int      `json:"pos"`
}

func (c *Client) GetTransaction(txHash string) (*Transaction, error) {
	rawRes, err := c.request("blockchain.transaction.get", []any{txHash})
	if err != nil {
		return nil, err
	}
	log.Printf("RECEIVED RAW TRANSACTION: %v", rawRes)

	var hexStr string
	if err := json.Unmarshal(rawRes, &hexStr); err != nil {
		return nil, err
	}
	txBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	btcutilTx, err := btcutil.NewTxFromBytes(txBytes)
	if err != nil {
		return nil, err
	}
	return btcutilToTransaction(btcutilTx), nil
}

func (c *Client) GetTransactionMerkle(txHash string, height int) (*TransactionMerkle, error) {
	rawTxMerkle, err := c.request("blockchain.transaction.get_merkle", []any{txHash, height})
	if err != nil {
		return nil, err
	}

	var txMerkle TransactionMerkle
	err = json.Unmarshal(rawTxMerkle, &txMerkle)
	if err != nil {
		return nil, err
	}
	return &txMerkle, nil
}

func btcutilToTransaction(utilTx *btcutil.Tx) *Transaction {
	msgTx := utilTx.MsgTx()

	tx := &Transaction{
		TxID: utilTx.Hash().String(),
		Vin: make([]struct {
			TxID string `json:"txid"`
			Vout int    `json:"vout"`
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
				Addresses []string `json:"addresses"`
			}{
				Addresses: addrStrings,
			},
		}
	}

	return tx
}
