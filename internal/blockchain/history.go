package blockchain

import (
	"log"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain/electrum"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"github.com/maphy9/btc-utxo-indexer/internal/util"
)

func (m *Manager) syncHistory(address string) error {
	txHdrs, err := m.client.GetTransactionHeaders(address)
	if err != nil {
		return err
	}

	for _, txHdr := range txHdrs {
		log.Printf("Received new transaction for address %s: %s", address, txHdr.TxHash)

		exists, err := m.db.Transactions().Exists(txHdr.TxHash)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		tx, err := m.client.GetTransaction(txHdr.TxHash)
		if err != nil {
			return err
		}

		txMerkle, err := m.client.GetTransactionMerkle(txHdr.TxHash, txHdr.Height)
		if err != nil {
			return err
		}

		hdr, err := m.db.Headers().GetByHeight(txHdr.Height)
		if err != nil {
			return err
		}
		if hdr == nil {
			continue // Header is out of sync
		}

		if !util.VerifyMerkleProof(txMerkle.Merkle, txHdr.TxHash, txMerkle.Pos, hdr.Root) {
			continue
		}

		// TODO: Put this in a transaction
		err = m.db.Transaction(func(q data.MasterQ) error {
			q.Transactions().Insert(TxHdrToData(txHdr))
	
			for _, in := range tx.Vin {
				err = q.Utxos().Spend(in.TxID, in.Vout, hdr.Height)
				if err != nil {
					return err
				}
			}
	
			for _, out := range tx.Vout {
				address := out.ScriptPubKey.Addresses[0]
				exists, err = q.Addresses().Exists(address)
				if err != nil {
					return err
				}
				if !exists {
					continue // Address is not tracked
				}
	
				_, err = q.Utxos().Insert(VoutToData(out, txHdr.TxHash, hdr.Height))
				if err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func VoutToData(vout electrum.UtxoVout, txHash string, height int) data.Utxo {
	sats := int64(vout.Value * 100_000_000)
	return data.Utxo{
		Address:       vout.ScriptPubKey.Addresses[0],
		TxHash:        txHash,
		TxPos:         vout.N,
		Value:         sats,
		CreatedHeight: height,
	}
}

func TxHdrToData(txHdr electrum.TransactionHeader) data.Transaction {
	return data.Transaction{
		Height: txHdr.Height,
		TxHash: txHdr.TxHash,
	}
}