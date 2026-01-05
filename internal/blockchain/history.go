package blockchain

import (
	"context"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain/electrum"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

func (m *Manager) processTransactionHeader(ctx context.Context, txHdr electrum.TransactionHeader) (*electrum.Transaction, error) {
	exists, err := m.db.Transactions().Exists(ctx, txHdr.TxHash)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, nil
	}

	tx, err := m.client.GetTransaction(ctx, txHdr.TxHash)
	if err != nil {
		return nil, err
	}

	txMerkle, err := m.client.GetTransactionMerkle(ctx, txHdr.TxHash, txHdr.Height)
	if err != nil {
		return nil, err
	}

	hdr, err := m.db.Headers().GetByHeight(ctx, txHdr.Height)
	if err != nil {
		return nil, err
	}
	if hdr == nil {
		return nil, nil
	}

	if !verifyMerkleProof(txMerkle.Merkle, txHdr.TxHash, txMerkle.Pos, hdr.Root) {
		return nil, nil
	}

	return tx, nil
}

func (m *Manager) saveTransaction(ctx context.Context, txHdr electrum.TransactionHeader, tx *electrum.Transaction) error {
	return m.db.Transaction(func(q data.MasterQ) error {
		_, err := q.Transactions().Insert(ctx, txHdrToData(txHdr))
		if err != nil {
			return err
		}

		for _, in := range tx.Vin {
			err = q.Utxos().Spend(in.TxID, in.Vout, tx.TxID)
			if err != nil {
				return err
			}
		}

		for _, out := range tx.Vout {
			address := out.ScriptPubKey.Addresses[0]
			exists, err := q.Addresses().Exists(address)
			if err != nil {
				return err
			}
			if !exists {
				continue
			}

			_, err = q.Utxos().Insert(voutToData(out, txHdr.TxHash, txHdr.Height))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (m *Manager) syncHistory(ctx context.Context, address string) error {
	txHdrs, err := m.client.GetTransactionHeaders(ctx, address)
	if err != nil {
		return err
	}

	for _, txHdr := range txHdrs {
		m.log.Infof("Received new transaction for address %s: %s", address, txHdr.TxHash)

		tx, err := m.processTransactionHeader(ctx, txHdr)
		if err != nil {
			return err
		}
		if tx == nil {
			continue
		}

		err = m.saveTransaction(ctx, txHdr, tx)
		if err != nil {
			return err
		}
	}

	return nil
}
