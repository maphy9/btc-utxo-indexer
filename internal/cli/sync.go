package cli

import (
	"context"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain"
	"github.com/maphy9/btc-utxo-indexer/internal/config"
	"github.com/maphy9/btc-utxo-indexer/internal/data/pg"
)

func SyncHeaders(cfg config.Config) error {
	db := pg.NewMasterQ(cfg.DB())
	log := cfg.Log()

	manager, err := blockchain.NewManager(cfg.ServiceConfig().NodeEntries, db, log)
	if err != nil {
		return err
	}
	defer func() {
		if err := manager.Close(); err != nil {
			log.WithError(err).Error("failed to close manager")
		}
	}()

	return manager.SyncHeaders(context.Background())
}
