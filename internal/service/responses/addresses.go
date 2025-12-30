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
