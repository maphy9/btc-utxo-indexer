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

	addrEntry := data.Address{
		UserID:  userID,
		Address: address,
	}

	return db.Transaction(func(q data.MasterQ) error {
		_, err := q.Addresses().Insert(ctx, addrEntry)
		if err != nil {
			return err
		}

		_, err = q.Utxos().InsertMany(ctx, utxos)
		if err != nil {
			return errors.New("utxos insert failed during address addition")
		}

		return nil
	})
}

func GetAddresses(ctx context.Context, db data.MasterQ, userID int64) ([]data.Address, error) {
	return db.Addresses().SelectByUserID(ctx, userID)
}

func CheckAddress(ctx context.Context, db data.MasterQ, userID int64, address string) (bool, error) {
	addr, err := db.Addresses().CheckAddress(ctx, userID, address)
	if err != nil {
		return false, err
	}
	return addr != nil, nil
}
