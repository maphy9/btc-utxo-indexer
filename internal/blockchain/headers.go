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
		rawHdrs, err := m.client.GetHeaders(height, chunkSize)
		if err != nil {
			return err
		}

		dataHdrs, err := headersToData(rawHdrs)
		if err != nil {
			return err
		}
		err = m.db.Headers().InsertBatch(dataHdrs)
		if err != nil {
			return err
		}
		log.Printf("Synchronized %d headers", height+len(rawHdrs))
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

func headersToData(rawHdrs []electrum.Header) ([]data.Header, error) {
	hdrs := make([]data.Header, len(rawHdrs))
	for i, rawHdr := range rawHdrs {
		hdr, err := util.ParseHeaderHex(rawHdr.Hex, rawHdr.Height)
		if err != nil {
			return nil, err
		}
		hdrs[i] = hdr
	}
	return hdrs, nil
}
