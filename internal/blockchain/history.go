package blockchain

import "github.com/maphy9/btc-utxo-indexer/internal/util"

func (m *Manager) syncHistory(address string) error {
	txHdrs, err := m.client.GetTransactionHeaders(address)
	if err != nil {
		return err
	}

	for _, txHdr := range txHdrs {
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
			return nil
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
		m.db.Transactions().Insert(txHdr.ToData())

		for _, in := range tx.Vin {
			err = m.db.Utxos().Spend(in.TxID, in.Vout, hdr.Height)
			if err != nil {
				return err
			}
		}

		for _, out := range tx.Vout {
			address := out.ScriptPubKey.Addresses[0]
			exists, err = m.db.Addresses().Exists(address)
			if err != nil {
				return err
			}
			if !exists {
				continue // Address is not tracked
			}

			_, err = m.db.Utxos().Insert(out.ToData(txHdr.TxHash, hdr.Height))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
