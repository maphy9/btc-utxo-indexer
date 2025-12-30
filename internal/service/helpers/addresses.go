package helpers

import (
	"net/http"

	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

func AddAddress(r *http.Request, addrStr string) error {
	ctx := r.Context()
	userID := UserID(r)
	db := DB(r)
	address := data.Address{
		UserID:  userID,
		Address: addrStr,
	}

	_, err := db.Addresses().Insert(ctx, address)
	return err
}

func GetAddresses(r *http.Request) ([]data.Address, error) {
	ctx := r.Context()
	userID := UserID(r)
	db := DB(r)
	return db.Addresses().SelectByUserID(ctx, userID)
}
