package requests

import (
	"encoding/json"
	"net/http"
)

type AddAddressRequest struct {
	Address string `json:"address"`
}

func NewAddAddressRequest(r *http.Request) (AddAddressRequest, error) {
	var request AddAddressRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	return request, err
}