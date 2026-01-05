package cli

import (
	"context"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain"
	"github.com/maphy9/btc-utxo-indexer/internal/blockchain/electrum"
	"github.com/maphy9/btc-utxo-indexer/internal/config"
	"github.com/maphy9/btc-utxo-indexer/internal/data/pg"
)

func SyncHeaders(cfg config.Config) error {
	db := pg.NewMasterQ(cfg.DB())
	log := cfg.Log()

	client, err := electrum.NewClient("electrum.blockstream.info:50002")
	if err != nil {
		return err
	}

	manager, err := blockchain.NewManager(client, db, log)
	if err != nil {
		return err
	}
	defer manager.Close()

	return manager.SyncHeaders(context.Background())
}
