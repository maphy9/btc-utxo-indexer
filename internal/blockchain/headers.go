package blockchain

import (
	"log"

	"github.com/maphy9/btc-utxo-indexer/internal/util"
)

const (
	chunkSize = 2016
)

func (m *Manager) SyncHeaders() error {
	localHeight, err := m.db.Headers().GetMaxHeight()
	if err != nil {
		return err
	}

	tipHeight, err := m.client.GetTipHeight()
	if err != nil {
		return err
	}

	for height := localHeight + 1; height <= tipHeight; height += chunkSize {
		hdrs, err := m.client.GetHeaders(height, chunkSize)
		if err != nil {
			return err
		}

		for _, hdr := range hdrs {
			dataHdr, err := util.ParseHeaderHex(hdr.Hex, hdr.Height)
			if err != nil {
				return err
			}
			_, err = m.db.Headers().Insert(dataHdr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Manager) ListenHeaders() error {
	notifyChan, err := m.client.SubscribeHeaders()
	if err != nil {
		return err
	}

	for hdr := range notifyChan {
		// Handle reorg
		log.Printf("New header found: %v", hdr)
	}

	return nil
}
