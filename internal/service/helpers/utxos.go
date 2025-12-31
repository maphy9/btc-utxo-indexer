package helpers

import (
	"context"

	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

func GetUtxos(ctx context.Context, db data.MasterQ, address string) ([]data.Utxo, error) {
	return db.Utxos().SelectByAddress(ctx, address)
}
