package cli

import (
	"github.com/maphy9/btc-utxo-indexer/internal/blockchain"
	"github.com/maphy9/btc-utxo-indexer/internal/config"
	"github.com/maphy9/btc-utxo-indexer/internal/data/pg"
)

func SyncHeaders(cfg config.Config) error {
	db := pg.NewMasterQ(cfg.DB())
	log := cfg.Log()

	manager, err := blockchain.NewManager("electrum.blockstream.info:50002", db, log)
	if err != nil {
		return err
	}
	return manager.SyncHeaders()
}
