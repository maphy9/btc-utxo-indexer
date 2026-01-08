package blockchain

import (
	"context"
	"sync"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain/electrum"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

func (m *Manager) processTransactionHeader(ctx context.Context, txHdr electrum.TransactionHeader) (*electrum.TransactionUtxos, error) {
	txMerkle, err := m.np.GetTransactionMerkle(ctx, txHdr.TxHash, txHdr.Height)
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

	return m.np.GetTransaction(ctx, txHdr.TxHash)
}

func (m *Manager) syncUtxos(ctx context.Context, createdUtxos []data.Utxo, spentUtxos []electrum.UtxoVin) error {
	return m.db.Transaction(func(q data.MasterQ) error {
		err := q.Utxos().InsertBatch(ctx, createdUtxos)
		if err != nil {
			return err
		}

		for _, spentUtxo := range spentUtxos {
			err = q.Utxos().Spend(ctx, spentUtxo.TxHash, spentUtxo.Vout, spentUtxo.SpentTxHash)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (m *Manager) syncTransactions(ctx context.Context, address string) error {
	txHdrs, err := m.np.GetTransactionHeaders(ctx, address)
	if err != nil {
		return err
	}

	err = m.db.Transactions().InsertBatch(ctx, txHdrsToData(txHdrs))
	if err != nil {
		return err
	}

	createdUtxos := make([]electrum.UtxoVout, 0, 64)
	spentUtxos := make([]electrum.UtxoVin, 0, 64)

	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	txHdrsChan := make(chan electrum.TransactionHeader, 10)
	errChan := make(chan error, 1)
	doneChan := make(chan struct{})
	once := sync.Once{}
	for i := 0; i < 10; i += 1 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for txHdr := range txHdrsChan {
				select {
				case <-doneChan:
					return
				default:
				}
				tx, err := m.processTransactionHeader(ctx, txHdr)
				if err != nil {
					once.Do(func() {
						errChan <- err
						close(doneChan)
					})
					return
				}
				if tx == nil {
					continue
				}

				mu.Lock()
				for _, utxo := range tx.Vouts {
					if utxo.Address != address {
						continue
					}
					createdUtxos = append(createdUtxos, utxo)
				}
				spentUtxos = append(spentUtxos, tx.Vins...)
				mu.Unlock()
			}
		}()
	}
	for _, txHdr := range txHdrs {
		txHdrsChan <- txHdr
	}
	close(txHdrsChan)
	wg.Wait()
	select {
	case err := <-errChan:
		return err
	default:
	}

	err = m.syncUtxos(ctx, voutsToData(createdUtxos), spentUtxos)
	if err != nil {
		return err
	}

	return nil
}
