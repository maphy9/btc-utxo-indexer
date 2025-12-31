package requests

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
)

type AddAddressRequest struct {
	Address string `json:"address"`
}

func NewAddAddressRequest(r *http.Request) (AddAddressRequest, error) {
	var request AddAddressRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return request, err
	}
	matched, err := regexp.Match("^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$", []byte(request.Address))
	if err != nil || !matched {
		return request, errors.New("invalid bitcoin address")
	}
	return request, nil
}
