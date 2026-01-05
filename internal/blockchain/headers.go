package blockchain

import (
	"log"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain/electrum"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"github.com/maphy9/btc-utxo-indexer/internal/util"
)

const (
	chunkSize = 2016
)

func (m *Manager) handleReorg(localTip, nextHdr *data.Header) (bool, error) {
	reorgDetected := false
	for nextHdr.ParentHash != localTip.Hash {
		reorgDetected = true

		err := m.db.Headers().DeleteByHeight(localTip.Height)
		if err != nil {
			return false, err
		}
		
		localTip, err = m.db.Headers().GetTipHeader()
		if err != nil {
			return false, err
		}
		if localTip == nil {
			break
		}

		rawNextHdr, err := m.client.GetHeader(localTip.Height + 1)
		if err != nil {
			return false, err
		}
		nextHdr, err = headerToData(rawNextHdr)
		if err != nil {
			return false, err
		}
	}
	return reorgDetected, nil
}

func (m *Manager) SyncHeaders() error {
	for {
		localTip, err := m.db.Headers().GetTipHeader()
		if err != nil {
			return err
		}
	
		tipHeight, err := m.client.GetTipHeight()
		if err != nil {
			return err
		}

		if localTip.Height >= tipHeight {
			break
		}

		startHeight := localTip.Height + 1
		count := min(chunkSize, tipHeight - startHeight + 1)
		rawHdrs, err := m.client.GetHeaders(startHeight, count)
		if err != nil {
			return err
		}

		dataHdrs, err := headersToData(rawHdrs)
		if err != nil {
			return err
		}

		if localTip.Height >= 0 {
			nextHdr := dataHdrs[0]
			reorgDetected, err := m.handleReorg(localTip, nextHdr)
			if err != nil {
				return err
			}
			if reorgDetected {
				continue
			}
		}

		err = m.db.Headers().InsertBatch(dataHdrs)
		if err != nil {
			return err
		}

		log.Printf("Synchronized headers %d-%d", startHeight, startHeight + count)
	}
	return nil
}

func (m *Manager) ListenHeaders() error {
	notifyChan, err := m.client.SubscribeHeaders()
	if err != nil {
		return err
	}

	for rawNextHdr := range notifyChan {
		log.Printf("Received header at height %d", rawNextHdr.Height)

		localTip, err := m.db.Headers().GetTipHeader()
		if err != nil {
			return err
		}
		if rawNextHdr.Height <= localTip.Height {
			continue
		}

		if rawNextHdr.Height > localTip.Height + 1 {
			if err := m.SyncHeaders(); err != nil {
				return err
			}
			continue
		}

		nextHdr, err := headerToData(&rawNextHdr)
		if err != nil {
			return err
		}
		reorgDetected, err := m.handleReorg(localTip, nextHdr)
		if reorgDetected {
			if err := m.SyncHeaders(); err != nil {
				return err
			}
		} else {
			_, err = m.db.Headers().Insert(nextHdr)
			if err != nil {
					return err
			}
		}
	}

	return nil
}

func headerToData(rawHdr *electrum.Header) (*data.Header, error) {
	return util.ParseHeaderHex(rawHdr.Hex, rawHdr.Height)
}

func headersToData(rawHdrs []electrum.Header) ([]*data.Header, error) {
	hdrs := make([]*data.Header, len(rawHdrs))
	for i, rawHdr := range rawHdrs {
		hdr, err := headerToData(&rawHdr)
		if err != nil {
			return nil, err
		}
		hdrs[i] = hdr
	}
	return hdrs, nil
}
