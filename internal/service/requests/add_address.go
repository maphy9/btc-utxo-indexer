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

var btcAddressPattern = regexp.MustCompile(`^([13][a-km-zA-HJ-NP-Z1-9]{25,34}|bc1[a-zA-HJ-NP-Z0-9]{11,71})$`)

func NewAddAddressRequest(r *http.Request) (AddAddressRequest, error) {
	var request AddAddressRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return request, err
	}
	if !btcAddressPattern.Match([]byte(request.Address)) {
		return request, errors.New("invalid bitcoin address")
	}
	return request, nil
}
