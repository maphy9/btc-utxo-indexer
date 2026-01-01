package blockchain

func (m *Manager) synchronizeHistory(address string) error {
	txHdrs, err := m.client.GetTransactionHeaders(address)
	if err != nil {
		return err
	}

	for _, txHdr := range txHdrs {
		exists, err := m.db.Transactions().Exists(txHdr.TxHash)
		if err != nil {
			return err
		}
		if !exists {
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

		block, err := m.db.Blocks().GetByHeight(txHdr.Height)
		if err != nil {
			return err
		}

		if !util.VerifyMerkleProof(txMerkle.Merkle, txHdr.TxHash, txMerkle.Pos, block.Root) {
			continue
		}

		// TODO: Put this in a transaction
		m.db.Transactions().Insert(tx.ToData())

		for _, in := range tx.Vin {
			err = m.db.Utxos().Spend(in.TxID, in.Vout)
			if err != nil {
				return err
			}
		}

		for _, out := range tx.Vout {
			address := out.ScriptPubKey.Addresses[0]
			if !m.db.Addresses().Exists(address) {
				continue
			}

			err = m.db.Utxos().Insert(out.ToData())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
