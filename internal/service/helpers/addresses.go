package helpers

import (
	"context"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

func AddAddress(ctx context.Context, db data.MasterQ, manager *blockchain.Manager, userID int64, address string) error {
	err := db.Transaction(func(q data.MasterQ) error {
		addressEntry, err := q.Addresses().InsertAddress(ctx, address)
		if err != nil {
			return err
		}

		_, err = q.Addresses().InsertUserAddress(ctx, data.UserAddress{
			AddressID: addressEntry.ID,
			UserID:    userID,
		})

		return err
	})
	if err != nil {
		return nil
	}

	return manager.WatchAddress(address)
}

func GetAddresses(ctx context.Context, db data.MasterQ, userID int64) ([]data.Address, error) {
	return db.Addresses().GetUserAddresses(ctx, userID)
}

func CheckAddress(ctx context.Context, db data.MasterQ, userID int64, address string) (bool, error) {
	userAddress, err := db.Addresses().GetUserAddress(ctx, userID, address)
	if err != nil {
		return false, err
	}
	return userAddress != nil, nil
}
