package responses

import "github.com/maphy9/btc-utxo-indexer/internal/data"

type GetAddressesResponse struct {
	Addresses []string `json:"addresses"`
}

func NewGetAddressesResponse(addresses []data.Address) GetAddressesResponse {
	response := make([]string, 0, len(addresses))
	for _, address := range addresses {
		response = append(response, address.Address)
	}
	return GetAddressesResponse{
		Addresses: response,
	}
}

type GetUtxosResponse struct {
	Utxos []data.Utxo `json:"utxos"`
}

func NewGetUtxosResponse(utxos []data.Utxo) GetUtxosResponse {
	return GetUtxosResponse{utxos}
}

type GetBalanceResponse struct {
	BalanceSat int64   `json:"balance_sat"`
	BalanceBtc float64 `json:"balance_btc"`
}

func NewGetBalanceResponse(balance int64) GetBalanceResponse {
	return GetBalanceResponse{
		BalanceSat: balance,
		BalanceBtc: float64(balance) / float64(100_000_000),
	}
}

type GetTransactionsResponse struct {
	Transactions []data.AddressTransaction `json:"transactions"`
}

func NewGetTransactionsResponse(transactions []data.AddressTransaction) GetTransactionsResponse {
	return GetTransactionsResponse{
		Transactions: transactions,
	}
}
