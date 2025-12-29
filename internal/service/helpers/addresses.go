package helpers

import (
	"net/http"

	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

func AddAddress(r *http.Request, userID int64, addr string) error {
	ctx := r.Context()
	db := DB(r)
	address := data.Address{
		UserID: userID,
		Address: addr,
	}

	_, err := db.Addresses().Insert(ctx, address)
	return err
}