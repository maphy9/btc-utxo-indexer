package blockchain

import (
	"encoding/hex"
	"errors"

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
		nextHdr, err = electrumHeaderToData(rawNextHdr)
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
		count := min(chunkSize, tipHeight-startHeight+1)
		rawHdrs, err := m.client.GetHeaders(startHeight, count)
		if err != nil {
			return err
		}

		dataHdrs, err := electrumHeadersToData(rawHdrs)
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

		m.log.Infof("synchronized headers %d-%d", startHeight, startHeight+count)
	}
	return nil
}

func (m *Manager) ListenHeaders() {
	notifyChan, err := m.client.SubscribeHeaders()
	if err != nil {
		m.log.WithError(err).Error("failed to subscribe to headers")
		return
	}

	for rawNextHdr := range notifyChan {
		m.log.Infof("received header at height %d", rawNextHdr.Height)

		localTip, err := m.db.Headers().GetTipHeader()
		if err != nil {
			m.log.WithError(err).Fatal("failed to get local tip header")
			continue
		}
		if rawNextHdr.Height <= localTip.Height {
			continue
		}

		if rawNextHdr.Height > localTip.Height+1 {
			if err := m.SyncHeaders(); err != nil {
				m.log.WithError(err).Error("failed to sync headers")
			}
			continue
		}

		nextHdr, err := electrumHeaderToData(&rawNextHdr)
		if err != nil {
			m.log.WithError(err).Fatal("failed to convert electrum header into data header")
			continue
		}
		reorgDetected, err := m.handleReorg(localTip, nextHdr)
		if reorgDetected {
			if err := m.SyncHeaders(); err != nil {
				m.log.WithError(err).Error("failed to sync headers")
				continue
			}
		} else {
			_, err = m.db.Headers().Insert(nextHdr)
			if err != nil {
				m.log.WithError(err).Fatal("failed to insert new header")
				continue
			}
		}
	}
}

func electrumHeaderToData(hdr *electrum.Header) (*data.Header, error) {
	if len(hdr.Hex) != 160 {
		return nil, errors.New("bad header hex")
	}
	bytes, err := hex.DecodeString(hdr.Hex)
	if err != nil {
		return nil, err
	}
	hash := hex.EncodeToString(util.Reverse(util.DoubleHash(bytes)))
	parentHash := hex.EncodeToString(util.Reverse(bytes[4:36]))
	root := hex.EncodeToString(util.Reverse(bytes[36:68]))
	return &data.Header{
		Hash:       hash,
		ParentHash: parentHash,
		Root:       root,
		Height:     hdr.Height,
	}, nil
}

func electrumHeadersToData(rawHdrs []electrum.Header) ([]*data.Header, error) {
	hdrs := make([]*data.Header, len(rawHdrs))
	for i, rawHdr := range rawHdrs {
		hdr, err := electrumHeaderToData(&rawHdr)
		if err != nil {
			return nil, err
		}
		hdrs[i] = hdr
	}
	return hdrs, nil
}
