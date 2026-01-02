package blockchain

import "log"

const (
	chunkSize = 2016
)

func (m *Manager) SyncHeaders() error {
	localHeight, err := m.db.Headers().GetMaxHeight()
	if err != nil {
		return err
	}

	tipHeader, err := m.client.GetTipHeader()
	if err != nil {
		return err
	}

	for height := localHeight + 1; height <= tipHeader.Height; height += chunkSize {
		rawHdrs, err := m.client.GetHeaders(height, chunkSize)
		if err != nil {
			return err
		}

		hdrs := rawHdrs.ToData()
		for _, hdr := range hdrs {
			_, err = m.db.Headers().Insert(hdr)
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
