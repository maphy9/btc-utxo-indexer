package helpers

import (
	"context"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

func AddAddress(ctx context.Context, db data.MasterQ, manager *blockchain.Manager, userID int64, address string) error {
	addr := data.Address{
		UserID:  userID,
		Address: address,
	}

	_, err := db.Addresses().Insert(ctx, addr)
	if err != nil {
		return err
	}

	utxos, err := manager.PrimaryNode.GetAddressUtxos(address)
	if err != nil || utxos == nil {
		return err
	}

	_, err = db.Utxos().InsertMany(ctx, utxos)
	return err
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
