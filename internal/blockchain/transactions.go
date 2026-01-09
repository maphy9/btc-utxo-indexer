package blockchain

import (
	"context"
	"sync"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain/rpc"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

func (m *Manager) processTransactionHeader(ctx context.Context, txHdr data.Transaction) (*rpc.TransactionData, error) {
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

	return m.np.GetTransactionData(ctx, txHdr.TxHash)
}

func (m *Manager) syncUtxos(ctx context.Context, txOuts []data.TransactionOutput, txIns []data.TransactionInput) error {
	return m.db.Transaction(func(q data.MasterQ) error {
		err := q.Transactions().InsertTransactionOutputsBatch(ctx, txOuts)
		if err != nil {
			return err
		}

		err = q.Transactions().InsertTransactionInputsBatch(ctx, txIns)
		if err != nil {
			return err
		}

		return q.Transactions().SpendTransactionOutputs(ctx, txIns)
	})
}

func (m *Manager) syncTransactions(ctx context.Context, address string) error {
	txHdrs, err := m.np.GetTransactionHeaders(ctx, address)
	if err != nil {
		return err
	}

	err = m.db.Transactions().InsertTransactionsBatch(ctx, txHdrs)
	if err != nil {
		return err
	}

	createdUtxos := make([]data.TransactionOutput, 0, 64)
	spentUtxos := make([]data.TransactionInput, 0, 64)

	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	healthyCount := m.np.GetHealthyCount()
	txHdrsChan := make(chan data.Transaction, healthyCount)
	var processingErr error
	doneChan := make(chan struct{})
	once := sync.Once{}
	for i := 0; i < healthyCount; i += 1 {
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
						processingErr = err
						close(doneChan)
					})
					return
				}
				if tx == nil {
					continue
				}

				mu.Lock()
				createdUtxos = append(createdUtxos, tx.Outputs...)
				spentUtxos = append(spentUtxos, tx.Inputs...)
				mu.Unlock()
			}
		}()
	}
	for _, txHdr := range txHdrs {
		txHdrsChan <- txHdr
	}
	close(txHdrsChan)
	wg.Wait()
	if processingErr != nil {
		return processingErr
	}

	err = m.syncUtxos(ctx, createdUtxos, spentUtxos)
	if err != nil {
		return err
	}

	return nil
}
