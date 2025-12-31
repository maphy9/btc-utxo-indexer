package helpers

import (
	"context"
	"net/http"

	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

func AddAddress(r *http.Request, address string) error {
	ctx := r.Context()
	userID := UserID(r)
	db := DB(r)
	manager := Manager(r)
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

	mappedUtxos := make([]data.Utxo, len(utxos), len(utxos))
	for i, utxo := range utxos {
		mappedUtxos[i] = data.Utxo{
			Address: address,
			TxID: utxo.TxID,
			Vout: utxo.Vout,
			Value: utxo.Value,
			BlockHeight: utxo.BlockHeight,
			BlockHash: utxo.BlockHash,
		}
	}

	_, err = db.Utxos().InsertMany(ctx, mappedUtxos)
	return err
}

func GetAddresses(r *http.Request) ([]data.Address, error) {
	ctx := r.Context()
	userID := UserID(r)
	db := DB(r)
	return db.Addresses().SelectByUserID(ctx, userID)
}

func CheckAddress(ctx context.Context, db data.MasterQ, userID int64, address string) (bool, error) {
	addr, err := db.Addresses().CheckAddress(ctx, userID, address)
	if err != nil {
		return false, err
	}
	return addr != nil, nil
}