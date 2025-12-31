package helpers

import (
	"context"
	"errors"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

func AddAddress(ctx context.Context, db data.MasterQ, manager *blockchain.Manager, userID int64, address string) error {
	utxos, err := manager.PrimaryNode.GetAddressUtxos(address)
	if err != nil || utxos == nil {
		return err
	}

	return db.Transaction(func(q data.MasterQ) error {
		addressEntry, err := q.Addresses().InsertAddress(ctx, address)
		if err != nil {
			return err
		}

		_, err = q.Addresses().InsertUserAddress(ctx, data.UserAddress{
			AddressID: addressEntry.ID,
			UserID:    userID,
		})
		if err != nil {
			return err
		}

		mappedUtxos := blockchain.MapRawUtxos(utxos, addressEntry.ID)
		_, err = q.Utxos().InsertMany(ctx, mappedUtxos)
		if err != nil {
			return errors.New("utxos insert failed during address addition")
		}

		return nil
	})
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
