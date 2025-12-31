package blockchain

import "github.com/maphy9/btc-utxo-indexer/internal/data"

type RawUtxo struct {
	TxID        string
	Vout        uint
	Value       int64
	BlockHeight int
	BlockHash   string
}

func (u *RawUtxo) MapRawUtxo(addressID int64) data.Utxo {
	return data.Utxo{
		AddressID:   addressID,
		TxID:        u.TxID,
		Vout:        u.Vout,
		Value:       u.Value,
		BlockHeight: u.BlockHeight,
		BlockHash:   u.BlockHash,
	}
}

func MapRawUtxos(rawUtxos []RawUtxo, addressID int64) []data.Utxo {
	result := make([]data.Utxo, 0, len(rawUtxos))
	for _, utxo := range rawUtxos {
		result = append(result, utxo.MapRawUtxo(addressID))
	}
	return result
}
