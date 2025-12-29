package responses

import "github.com/maphy9/btc-utxo-indexer/internal/data"

type GetAddressesResponse struct {
	Addresses []string `json:"addresses"`
}

func NewGetAddressesResponse(addresses []data.Address) GetAddressesResponse {
	response := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		response = append(response, addr.Address)
	}
	return GetAddressesResponse{
		Addresses: response,
	}
}
