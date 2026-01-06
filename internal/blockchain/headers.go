package blockchain

import (
	"context"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain/electrum"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

const (
	chunkSize = 2016
)

func (m *Manager) handleReorg(ctx context.Context, localTip, nextHdr *data.Header) (bool, error) {
	reorgDetected := false
	for nextHdr.ParentHash != localTip.Hash {
		reorgDetected = true

		err := m.db.Headers().DeleteByHeight(ctx, localTip.Height)
		if err != nil {
			return false, err
		}

		localTip, err = m.db.Headers().GetTipHeader(ctx)
		if err != nil {
			return false, err
		}
		if localTip == nil {
			localTip = &data.Header{
				Height: -1,
			}
		}

		rawNextHdr, err := m.client.GetHeader(ctx, localTip.Height+1)
		if err != nil {
			return false, err
		}
		nextHdr, err = electrumHeaderToData(rawNextHdr)
		if err != nil {
			return false, err
		}
	}
	return reorgDetected, nil
}

func (m *Manager) SyncHeaders(ctx context.Context) error {
	for {
		localTip, err := m.db.Headers().GetTipHeader(ctx)
		if err != nil {
			return err
		}
		if localTip == nil {
			localTip = &data.Header{
				Height: -1,
			}
		}

		tipHeight, err := m.client.GetTipHeight(ctx)
		if err != nil {
			return err
		}

		if localTip.Height >= tipHeight {
			break
		}

		startHeight := localTip.Height + 1
		count := min(chunkSize, tipHeight-startHeight+1)
		rawHdrs, err := m.client.GetHeaders(ctx, startHeight, count)
		if err != nil {
			return err
		}

		dataHdrs, err := electrumHeadersToData(rawHdrs)
		if err != nil {
			return err
		}

		if localTip.Height >= 0 {
			nextHdr := dataHdrs[0]
			reorgDetected, err := m.handleReorg(ctx, localTip, nextHdr)
			if err != nil {
				return err
			}
			if reorgDetected {
				continue
			}
		}

		err = m.db.Headers().InsertBatch(ctx, dataHdrs)
		if err != nil {
			return err
		}

		m.log.Infof("synchronized headers %d-%d", startHeight, startHeight+count)
	}
	return nil
}

func (m *Manager) ListenHeaders() {
	m.wg.Add(1)
	defer m.wg.Done()

	notifyChan, err := m.client.SubscribeHeaders(m.ctx)
	if err != nil {
		m.log.WithError(err).Error("failed to subscribe to headers")
		return
	}

	for {
		select {
		case <-m.ctx.Done():
			return
		case rawNextHdr, ok := <-notifyChan:
			if !ok {
				m.log.Info("Headers channed was closed")
				return
			}
			m.log.Infof("received header at height %d", rawNextHdr.Height)
			err = m.processHeader(m.ctx, rawNextHdr)
			if err != nil {
				m.log.WithError(err).Fatal("failed to process header")
				continue
			}
		}
	}
}

func (m *Manager) processHeader(ctx context.Context, rawNextHdr electrum.Header) error {
	localTip, err := m.db.Headers().GetTipHeader(ctx)
	if err != nil {
		return err
	}
	if localTip == nil {
		localTip = &data.Header{
			Height: -1,
		}
	}

	if rawNextHdr.Height <= localTip.Height {
		return err
	}

	if rawNextHdr.Height > localTip.Height+1 {
		if err := m.SyncHeaders(ctx); err != nil {
			return err
		}
		return nil
	}

	nextHdr, err := electrumHeaderToData(&rawNextHdr)
	if err != nil {
		return err
	}

	reorgDetected, err := m.handleReorg(ctx, localTip, nextHdr)
	if err != nil {
		return err
	}
	if reorgDetected {
		if err := m.SyncHeaders(ctx); err != nil {
			return err
		}
		return nil
	}

	_, err = m.db.Headers().Insert(ctx, nextHdr)
	return err
}
